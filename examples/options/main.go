package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/options"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	Options  options.Model
	Selected bool
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		Options: options.Model{
			Options: []string{
				"Tobi",
				"Loki",
				"Jane",
				"Manny",
				"Luna",
			},
		},
	}, nil
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)

	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyEnter:
			if m.Selected {
				return m, tea.Quit
			}
			m.Selected = true
			return m, nil
		case terminput.KeyEscape:
			return m, tea.Quit
		case terminput.KeyRune:
			switch r := msg.Rune(); r {
			case 'q':
				return m, tea.Quit
			default:
				if !m.Selected {
					m.Options = options.Update(msg, m.Options)
				}
				return m, nil
			}
		default:
			if !m.Selected {
				m.Options = options.Update(msg, m.Options)
			}
			return m, nil
		}
	}

	return m, nil
}

// view function.
func view(ctx context.Context, model tea.Model) string {
	w := new(bytes.Buffer)
	m := model.(Model)

	// padding
	fmt.Fprintf(w, "\r\n")
	defer fmt.Fprintf(w, "\r\n")

	// input
	if m.Selected {
		fmt.Fprintf(w, "  You chose:\r\n")
		for _, o := range m.Options.Value() {
			fmt.Fprintf(w, "  - %s\r\n", o)
		}
	} else {
		fmt.Fprintf(w, "  Choose your favorite pet:\r\n\r\n")
		fmt.Fprintf(w, "%s", options.View(m.Options))
	}

	return w.String()
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\r\n", err)
	}
}
