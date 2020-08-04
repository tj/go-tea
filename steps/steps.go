// Package steps provides a wizard style step progress bar.
package steps

import "strings"

// stepCompletedChar is the character used for a completed step.
var stepCompletedChar = "◉"

// stepChar is the character used for a step.
var stepChar = "◯"

// barCompletedChar is the character used for the completed portion of the bar.
var barCompletedChar = "━"

// barChar is the character used for the uncompleted portion of the bar.
var barChar = "━"

// Model is the step model.
type Model struct {
	// Steps is the set of step labels.
	Steps []string

	// Step is the index of the current step.
	Step int

	// Prefix is a string applied to the beginning of each line.
	Prefix string
}

// View function.
func View(m Model) (s string) {
	pad := 8
	max := maxLength(m.Steps) + pad

	s += m.Prefix

	// progress bar
	for i := range m.Steps {
		complete := i < m.Step

		if complete {
			s += stepCompletedChar
		} else {
			s += stepChar
		}

		if i < len(m.Steps)-1 {
			if complete {
				s += strings.Repeat(barCompletedChar, max)
			} else {
				s += strings.Repeat(barChar, max)
			}
		}
	}

	s += "\n" + m.Prefix

	// step labels
	var lpad int
	for i, step := range m.Steps {
		switch {
		// first step, left align
		case i == 0:
			s += step
			lpad = max + 1 - len(step)
		// last step, right align
		case i == len(m.Steps)-1:
			lpad -= len(step)
			s += strings.Repeat(" ", lpad) + step
		// others, center
		default:
			lpad -= len(step) / 2
			s += strings.Repeat(" ", lpad) + step
			lpad = max + 1 - len(step)/2
		}
	}

	s += "\n"

	return
}

// maxLength returns the max length of the given strings.
func maxLength(values []string) (max int) {
	for _, s := range values {
		if len(s) > max {
			max = len(s)
		}
	}
	return
}
