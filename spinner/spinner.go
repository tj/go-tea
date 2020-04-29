// Package spinner provides a frame-based textual spinner.
package spinner

import (
	"context"
	"time"

	"github.com/tj/go-tea"
)

// Tick message.
type Tick struct{}

// DefaultInterval is the default animation interval used.
var DefaultInterval = time.Millisecond * 75

// DefaultFrames is the default set of frames used.
var DefaultFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

// Model is the input model.
type Model struct {
	// Frames is a set of frames to animation. Defaults to DefaultFrames.
	Frames []string

	// Interval is the animation update interval. Defaults to DefaultInterval.
	Interval time.Duration

	// current tick.
	tick int
}

// Update function.
func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	interval := m.Interval

	if interval == 0 {
		interval = DefaultInterval
	}

	if _, ok := msg.(Tick); ok {
		m.tick++
		return m, tick(interval)
	}

	return m, nil
}

// View function.
func View(m Model) string {
	frames := m.Frames
	if frames == nil {
		frames = DefaultFrames
	}
	frame := (m.tick + 1) % len(frames)
	return frames[frame]
}

// tick is a command which advances the spinner animation frame.
func tick(d time.Duration) tea.Cmd {
	return func(ctx context.Context) tea.Msg {
		time.Sleep(d)
		return Tick{}
	}
}
