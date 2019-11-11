package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/input"
	"github.com/tj/go-terminput"

	"github.com/tj/go-tea/examples/todo/list"
)

// AddedItem msg.
type AddedItem string

// Model struct.
type Model struct {
	// List is the todo items.
	List list.Model

	// Input is the text input used for a new item.
	Input input.Model

	// FocusingAddItem is used to indicate the "Add Item"
	// button is in focus.
	FocusingAddItem bool

	// AddingItem is used to indicate that we're adding
	// a new item.
	AddingItem bool
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		List: list.Model{
			Items: []string{
				"Buy groceries",
				"Feed the ferrets",
				"Feed the cats",
			},
		},
	}, nil
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)

	// delegate messages to input
	if m.AddingItem {
		m.Input = input.Update(msg, m.Input)
	}

	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyUp:
			if m.FocusingAddItem {
				m.List.Disabled = false
				m.FocusingAddItem = false
			} else {
				m.List = list.Update(msg, m.List)
			}
			return m, nil
		case terminput.KeyDown:
			// we were already at the end of the list
			if m.List.Selected == len(m.List.Items)-1 {
				// focus the add "button"
				m.FocusingAddItem = true
				m.List.Disabled = true
				return m, nil
			}

			m.List.Disabled = false
			m.List = list.Update(msg, m.List)
			return m, nil
		case terminput.KeyEnter:
			// add a new item, clear the input, select the last one
			if m.AddingItem {
				m.List.Items = append(m.List.Items, m.Input.Value)
				m.Input = input.Model{}
				m.AddingItem = false
				m.List.Selected = len(m.List.Items) - 1
				return m, nil
			}

			// we're on the add button, show the input
			if m.FocusingAddItem {
				m.AddingItem = true
				return m, nil
			}

			return m, nil
		case terminput.KeyBackspace:
			m.List = list.Update(msg, m.List)
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
	fmt.Fprintf(w, "\r\n")
	defer fmt.Fprintf(w, "\r\n")

	// help
	if !m.AddingItem {
		fmt.Fprintf(w, "  Use the arrows to navigate, backspace to remove, and 'q' to exit.\r\n\r\n")
	}

	// adding
	if m.AddingItem {
		fmt.Fprintf(w, "  Item (press enter): %s", input.View(m.Input))
		return w.String()
	}

	// listing
	fmt.Fprintf(w, list.View(m.List))

	// add button
	if m.FocusingAddItem {
		fmt.Fprintf(w, "\r\n  \033[1mAdd Item\033[m\r\n")
	} else {
		fmt.Fprintf(w, "\r\n  Add Item\r\n")
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
