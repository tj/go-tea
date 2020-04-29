package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/spinner"
	"github.com/tj/go-terminput"
)

// start command.
func start(ctx context.Context) tea.Msg {
	return nil
}

// Model struct.
type Model struct {
	Spinner spinner.Model
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	// the start command is only necessary to trigger an update() below,
	// which in turn allows spinner.Update() to advance the frame. The
	// start command itself doesn't do anything.
	return Model{}, start
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)

	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		// pressed esc or q
		if msg.Key() == terminput.KeyEscape || msg.Rune() == 'q' {
			return m, tea.Quit
		}
	}

	spinner, cmd := spinner.Update(msg, m.Spinner)
	m.Spinner = spinner
	return m, cmd
}

// view function.
func view(ctx context.Context, model tea.Model) string {
	m := model.(Model)
	return fmt.Sprintf("\n  Deploying %s\n\n  [q] Quit", green(spinner.View(m.Spinner)))
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\r\n", err)
	}
}

// green string.
func green(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}
