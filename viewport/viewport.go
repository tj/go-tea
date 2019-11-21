// Package viewport provides a scrollable view into content.
package viewport

import (
	"fmt"
	"strings"

	"github.com/tj/go-tea"
	"github.com/tj/go-terminput"
)

// Model is the viewport model.
type Model struct {
	// Height is the viewport height, usually the terminal height.
	Height int

	// ScrollHeight is the scrollable height.
	ScrollHeight int

	// ScrollY is the vertical scroll position.
	ScrollY int

	// ScrollBy is the number of rows or columns to scroll by.
	ScrollBy int
}

// Update function.
func Update(msg tea.Msg, m Model) Model {
	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyUp:
			m.ScrollY = max(0, m.ScrollY-m.ScrollBy)
			return m
		case terminput.KeyDown:
			m.ScrollY = min(m.ScrollY+m.ScrollBy, m.ScrollHeight-m.Height)
			return m
		}
	}
	return m
}

// View function.
func View(m Model, content string) string {
	lines := strings.Split(content, "\n")
	from := m.ScrollY
	to := m.ScrollY + m.Height
	lines = bounded(lines, from, to)
	return strings.Join(lines, "\n")
}

// bounded slice.
func bounded(s []string, from, to int) []string {
	from = max(0, min(from, len(s)))
	to = max(0, min(to, len(s)))
	return s[from:to]
}

// min returns the minimum of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// bell sound.
func bell() {
	fmt.Printf("\a")
}
