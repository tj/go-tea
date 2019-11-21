package main

import (
	"context"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/shortcut"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	Message   string
	Shortcuts shortcut.Model
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		Message: "Hello World",
		Shortcuts: shortcut.Model{
			Keys: []shortcut.Key{
				{
					Key:   "q",
					Label: "Quit",
				},
			},
		},
	}, nil
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
		return m, nil
	}

	return m, nil
}

// view function.
func view(ctx context.Context, model tea.Model) string {
	m := model.(Model)
	return "\n" + m.Message + "\n\n" + shortcut.View(m.Shortcuts)
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
