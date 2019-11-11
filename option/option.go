package option

import (
	"bytes"
	"fmt"

	"github.com/tj/go-tea"
	"github.com/tj/go-terminput"
)

// Model is the option input model.
type Model struct {
	// Options is the set of options the user can select.
	Options []string

	// Selected is the index of the selected value.
	Selected int
}

// Value returns the selected option.
func (m *Model) Value() string {
	return m.Options[m.Selected]
}

// Update function.
func Update(msg tea.Msg, m Model) Model {
	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyUp:
			if m.Selected > 0 {
				m.Selected--
			} else {
				bell()
			}
		case terminput.KeyDown:
			if m.Selected < len(m.Options)-1 {
				m.Selected++
			} else {
				bell()
			}
		}
	}
	return m
}

// View function.
func View(m Model) string {
	w := new(bytes.Buffer)

	for i, option := range m.Options {
		if i == m.Selected {
			fmt.Fprintf(w, "  \033[1m%s\033[m\r\n", option)
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
