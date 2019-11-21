package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/option"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	Option   option.Model
	Selected bool
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		Option: option.Model{
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
					m.Option = option.Update(msg, m.Option)
				}
				return m, nil
			}
		default:
			if !m.Selected {
				m.Option = option.Update(msg, m.Option)
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
	fmt.Fprintf(w, "\n")
	defer fmt.Fprintf(w, "\n")

	// input
	if m.Selected {
		fmt.Fprintf(w, "  You chose: %s\n", m.Option.Value())
	} else {
		fmt.Fprintf(w, "  Choose your favorite pet:\n\n")
		fmt.Fprintf(w, "%s", option.View(m.Option))
	}

	return w.String()
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
