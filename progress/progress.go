// Package progress provides a progress bar.
package progress

import (
	"fmt"
	"math"
	"strings"
)

// Model is the progress bar model.
type Model struct {
	// Filled bar character, defaulting to '█'.
	Filled string

	// Empty bar character, defaulting to '░'.
	Empty string

	// Width of the progress bar, defaulting to 24.
	Width int

	// Percentage of completion from 0-1.
	Percent float64
}

// View function.
func View(m Model) string {
	w := m.Width
	if w == 0 {
		w = 24
	}
	f := defaultString(m.Filled, "█")
	e := defaultString(m.Empty, "░")
	nf := int(math.Ceil(float64(w) * m.Percent))
	ne := w - nf
	bar := strings.Repeat(f, nf) + strings.Repeat(e, ne)
	return fmt.Sprintf("|%s| %0.0f%%", bar, m.Percent*100)
}

// defaultString returns s or the default.
func defaultString(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
