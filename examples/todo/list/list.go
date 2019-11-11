package list

import (
	"bytes"
	"fmt"

	"github.com/tj/go-tea"
	"github.com/tj/go-terminput"
)

// Model is the todo list model.
type Model struct {
	// Items is the set of todo items.
	Items []string

	// Selected is the index of the selected item.
	Selected int

	// Disabled will disable the list when true.
	Disabled bool

	// Removing is true when an item is pending removal.
	Removing bool
}

// Value returns the selected option.
func (m *Model) Value() string {
	return m.Items[m.Selected]
}

// Update function.
func Update(msg tea.Msg, m Model) Model {
	if m.Disabled {
		return m
	}
	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyUp:
			if m.Selected > 0 {
				m.Selected--
			} else {
				bell()
			}
			m.Removing = false
		case terminput.KeyDown:
			if m.Selected < len(m.Items)-1 {
				m.Selected++
			} else {
				bell()
			}
			m.Removing = false
		case terminput.KeyBackspace:
			if m.Removing {
				m.Items = append(m.Items[:m.Selected], m.Items[m.Selected+1:]...)
				m.Removing = false
			} else {
				m.Removing = true
			}
		}
	}
	return m
}

// View function.
func View(m Model) string {
	w := new(bytes.Buffer)

	for i, option := range m.Items {
		if i == m.Selected && !m.Disabled {
			if m.Removing {
				fmt.Fprintf(w, "  \033[1;31m%s (press again to confirm removal)\033[m\r\n", option)
			} else {
				fmt.Fprintf(w, "  \033[1m%s\033[m\r\n", option)
			}
		} else {
			fmt.Fprintf(w, "  %s\r\n", option)
		}
	}

	return w.String()
}

// bell sound.
func bell() {
	fmt.Printf("\a")
}
