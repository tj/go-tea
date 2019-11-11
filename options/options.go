// Package options provides an options list with many selectable values.
package options

import (
	"bytes"
	"fmt"

	"github.com/tj/go-tea"
	"github.com/tj/go-terminput"
)

// Model is the options input model.
type Model struct {
	// Options is the set of options the user can select.
	Options []string

	// Selected is the indexes of the selected values.
	Selected []int

	// active index.
	index int
}

// Value returns the selected option.
func (m *Model) Value() (values []string) {
	for _, i := range m.Selected {
		if i < len(m.Options) {
			values = append(values, m.Options[i])
		}
	}
	return
}

// Update function.
func Update(msg tea.Msg, m Model) Model {
	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyUp:
			if m.index > 0 {
				m.index--
			} else {
				bell()
			}
		case terminput.KeyDown:
			if m.index < len(m.Options)-1 {
				m.index++
			} else {
				bell()
			}
		case terminput.KeyRune:
			if msg.Rune() == ' ' {
				return toggle(m)
			}
		}
	}
	return m
}

// View function.
func View(m Model) string {
	w := new(bytes.Buffer)

	for i, option := range m.Options {
		if i == m.index {
			fmt.Fprintf(w, "\033[1m")
		}

		if isSelected(m, i) {
			fmt.Fprintf(w, "  ■ %s\r\n", option)
		} else {
			fmt.Fprintf(w, "  □ %s\r\n", option)
		}

		if i == m.index {
			fmt.Fprintf(w, "\033[0m")
		}
	}

	return w.String()
}

// toggle selection at the current index.
func toggle(m Model) Model {
	if isSelected(m, m.index) {
		for i, v := range m.Selected {
			if v == m.index {
				m.Selected = append(m.Selected[:i], m.Selected[i+1:]...)
			}
		}
	} else {
		m.Selected = append(m.Selected, m.index)
	}
	return m
}

// isSelected returns true if the index is selected.
func isSelected(m Model, index int) bool {
	for _, i := range m.Selected {
		if i == index {
			return true
		}
	}
	return false
}

// bell sound.
func bell() {
	fmt.Printf("\a")
}
