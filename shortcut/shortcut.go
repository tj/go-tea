// Package shortcut provides a shortcut keys.
package shortcut

import (
	"fmt"
)

// Key is a shortcut key.
type Key struct {
	// Key is the key combination such as "q".
	Key string

	// Label is the key label such as "Quit".
	Label string
}

// Model is the input model.
type Model struct {
	// Keys is a set of keyboard shortcuts available.
	Keys []Key
}

// View function.
func View(m Model) (s string) {
	for _, k := range m.Keys {
		s += fmt.Sprintf("[%s] %s ", k.Key, k.Label)
	}
	return
}
