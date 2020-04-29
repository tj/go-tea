package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/input"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	Input   input.Model
	Editing bool
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		Input:   input.Model{},
		Editing: true,
	}, nil
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)
	m.Input = input.Update(msg, m.Input)

	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyEnter:
			m.Editing = false
			return m, tea.Quit
		case terminput.KeyEscape:
			return m, tea.Quit
		case terminput.KeyRune:
			switch r := msg.Rune(); r {
			case 'q':
				return m, tea.Quit
			}
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

	// editing
	if m.Editing {
		fmt.Fprintf(w, "  Enter your name: %s\n", input.View(m.Input))
	} else {
		fmt.Fprintf(w, "  Hello %s!\n\n", m.Input.Value)
	}

	return w.String()
}

func bell() {
	fmt.Printf("\a")
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
