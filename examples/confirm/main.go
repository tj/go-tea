package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/input"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	ProjectID string
	Confirm   input.Model
	Confirmed bool
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		ProjectID: "my-project",
	}, nil
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)
	m.Confirm = input.Update(msg, m.Confirm)

	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyEnter:
			if m.Confirm.Value == m.ProjectID {
				m.Confirmed = true
				return m, tea.Quit
			}
			bell()
			return m, nil
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

	// confirm
	if m.Confirmed {
		fmt.Fprintf(w, "  Deleted %q.\n\n", m.ProjectID)
	} else if !strings.HasPrefix(m.ProjectID, m.Confirm.Value) {
		fmt.Fprintf(w, "  Enter %q to confirm deletion: %s\n", m.ProjectID, red(input.View(m.Confirm)))
	} else {
		fmt.Fprintf(w, "  Enter %q to confirm deletion: %s\n", m.ProjectID, input.View(m.Confirm))
	}

	return w.String()
}

// red string.
func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

// bell sound.
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
