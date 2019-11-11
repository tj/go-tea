package main

import (
	"context"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	Message  string
	Confused bool
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		Message: "Hello World",
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

		// pressed some other key
		m.Confused = true
		return m, nil
	}

	return m, nil
}

// view function.
func view(ctx context.Context, model tea.Model) string {
	m := model.(Model)
	if m.Confused {
		return m.Message + " (press q to exit)"
	}
	return m.Message
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\r\n", err)
	}
}
