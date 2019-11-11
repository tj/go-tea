package input

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/tj/go-tea"
	"github.com/tj/go-terminput"
)

// Model is the input model.
type Model struct {
	// Value is the text input value.
	Value string

	// pos is the position of the cursor.
	pos int
}

// Update function.
func Update(msg tea.Msg, m Model) Model {
	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyBackspace:
			if m.pos > 0 {
				m.Value = m.Value[:m.pos-1] + m.Value[m.pos:]
				m.pos--
			} else {
				bell()
			}
			return m
		case terminput.KeyLeft:
			if m.pos > 0 {
				if msg.Alt() {
					m.pos -= wordLeft(m)
				} else {
					m.pos--
				}
			} else {
				bell()
			}
			return m
		case terminput.KeyRight:
			if m.pos < len(m.Value) {
				if msg.Alt() {
					m.pos += wordRight(m)
				} else {
					m.pos++
				}
			} else {
				bell()
			}
			return m
		case terminput.KeyRune:
			m.Value = m.Value[:m.pos] + string(msg.Rune()) + m.Value[m.pos:]
			m.pos++
			return m
		}
	}
	return m
}

// View function.
func View(m Model) string {
	if m.Value == "" {
		return cursor(" ")
	}
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "%s", m.Value[:m.pos])
	if m.pos < len(m.Value) {
		fmt.Fprintf(w, "")
		fmt.Fprintf(w, cursor(string(m.Value[m.pos])))
		fmt.Fprintf(w, m.Value[m.pos+1:])
	} else {
		fmt.Fprintf(w, cursor(" "))
	}
	return w.String()
}

// wordLeft util.
func wordLeft(m Model) (size int) {
	// TODO: support utf8
	i := m.pos - 1

	// skip whitespace
	for i >= 0 {
		if unicode.IsSpace(rune(m.Value[i])) {
			size++
			i--
		} else {
			break
		}
	}

	// skip word
	for i >= 0 {
		if !unicode.IsSpace(rune(m.Value[i])) {
			size++
			i--
		} else {
			break
		}
	}

	return
}

// wordRight util.
func wordRight(m Model) (size int) {
	// TODO: support utf8
	i := m.pos

	// skip word
	for i < len(m.Value) {
		if !unicode.IsSpace(rune(m.Value[i])) {
			size++
			i++
		} else {
			break
		}
	}

	// skip whitespace
	for i < len(m.Value) {
		if unicode.IsSpace(rune(m.Value[i])) {
			size++
			i++
		} else {
			break
		}
	}

	return
}

// isWord character.
func isWord(c rune) bool {
	return !unicode.IsSpace(c)
}

// cursor styling.
func cursor(s string) string {
	return fmt.Sprintf("\033[48;5;61m%s\033[0m", s)
}

// bell sound.
func bell() {
	fmt.Printf("\a")
}
