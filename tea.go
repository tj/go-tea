// Package tea provides an Elm inspired functional framework for interactive command-line programs.
package tea

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/term"
	"github.com/tj/go-terminput"
)

// quitMsg is the internal message for exiting the program.
type quitMsg struct{}

// batchMsg is the internal message for performing a batch of commands.
type batchMsg []Cmd

// Msg is passed to your program's Update() function, representing an
// action which was performed, for example a ItemRemoved msg might be
// a struct containing the ID of the item removed.
type Msg interface{}

// Model is the model which defines all or a subset of your program state.
type Model interface{}

// Init is a function which is invoked when starting your program, returning
// the initial model and optional command.
type Init func(context.Context) (Model, Cmd)

// Update is a function which is invoked for every message, allowing you
// to return a new, updated model and optional command.
type Update func(context.Context, Msg, Model) (Model, Cmd)

// Cmd is a function used to perform an action, when complete you may
// return a message, error, or nil.
//
// For example a comand which removes a user from a database might
// return a struct UserRemoved with its ID so that Update()
// can remove it before rendering.
//
// Errors are special cased, so you may return an error in place
// of a message, this will cause the program to exit and the
// error will be printed. If you wish to handle errors in
// a different way, you should return a message containing
// the error and update your model accordingly.
//
// Returning nil is a no-op.
//
type Cmd func(context.Context) Msg

// View is a function used to render the program's model
// before it is written to the terminal.
type View func(context.Context, Model) string

// Quit is a message which exits the program.
//
// For example:
//
//   return m, tea.Quit
//
func Quit(ctx context.Context) Msg {
	return quitMsg{}
}

// Batch performs many commands concurrently,
// with no order guarantees.
func Batch(cmds ...Cmd) Cmd {
	return func(ctx context.Context) Msg {
		return batchMsg(cmds)
	}
}

// Program is a terminal application comprised init,
// update, and view functions.
type Program struct {
	// Init function.
	Init

	// Update function.
	Update

	// View function.
	View

	rw io.ReadWriter
}

// NewProgram returns a new program.
func NewProgram(init Init, update Update, view View) *Program {
	return &Program{
		Init:   init,
		Update: update,
		View:   view,
	}
}

// Start the program.
func (p *Program) Start(ctx context.Context) error {
	// open tty
	tty, err := term.Open("/dev/tty")
	if err != nil {
		return err
	}
	p.rw = tty

	// raw mode
	tty.SetRaw()
	defer tty.Restore()

	// hide cursor
	hideCursor()
	defer showCursor()

	return p.start(ctx)
}

// start implementation.
func (p *Program) start(ctx context.Context) error {
	msgs := make(chan Msg)
	cmds := make(chan Cmd)
	done := make(chan struct{})
	errs := make(chan error)

	// input loop. We read user input and provide
	// them to the application as msgs.
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				msg, err := terminput.Read(p.rw)
				if err != nil {
					errs <- err
					return
				}
				msgs <- msg
			}
		}
	}()

	// command loop. We asynchronously process
	// any commands received in the background,
	// which may produce msgs.
	go func() {
		for {
			select {
			case <-done:
				return
			case cmd := <-cmds:
				if cmd != nil {
					go func() {
						msgs <- cmd(ctx)
					}()
				}
			}
		}
	}()

	// initialize app
	model, cmd := p.Init(ctx)
	cmds <- cmd

	// draw the initial view
	prev := normalize(p.View(ctx, model))
	io.WriteString(p.rw, prev)

	// draw loop. We process msgs, passing them
	// to the Update() function followed by the
	// View() function for rendering.
	for {
		select {
		case err := <-errs:
			close(done)
			return err
		case msg := <-msgs:
			// quit msg
			if _, ok := msg.(quitMsg); ok {
				close(done)
				return nil
			}

			// error msg
			if err, ok := msg.(error); ok {
				close(done)
				return err
			}

			// batch msg
			if v, ok := msg.(batchMsg); ok {
				for _, cmd := range v {
					cmds <- cmd
				}
				continue
			}

			// update
			model, cmd = p.Update(ctx, msg, model)
			cmds <- cmd

			// render view changes
			curr := normalize(p.View(ctx, model))
			clearLines(strings.Count(prev, "\r\n") + 1)
			io.WriteString(p.rw, curr)
			prev = curr
		}
	}
}

// normalize .
func normalize(s string) string {
	return strings.Replace(s, "\n", "\r\n", -1)
}

// hideCursor hides the cursor.
func hideCursor() {
	fmt.Printf("\033[?25l")
}

// showCursor shows the cursor.
func showCursor() {
	fmt.Printf("\033[?25h")
}

// clearLines clears a number of lines.
func clearLines(n int) {
	for i := 0; i < n; i++ {
		moveUp(1)
		clearLine()
	}
}

// clearLine clears the entire line.
func clearLine() {
	fmt.Printf("\033[2K")
}

// moveUp moves the cursor to the beginning of n lines up.
func moveUp(n int) {
	fmt.Printf("\033[%dF", n)
}

// clear the screen.
func clear() {
	fmt.Printf("\033[2J\033[3J\033[1;1H")
}
