package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/progress"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	URL           string
	Progress      progress.Model
	Iteration     int
	MaxIterations int
	Previous      struct {
		StatusCode int
		Duration   time.Duration
	}
	Durations []time.Duration
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	url := "https://apex.sh"
	return Model{
		URL:           url,
		Iteration:     0,
		MaxIterations: 500,
	}, request(url)
}

// requestCompleted message.
type requestCompleted struct {
	StatusCode int
	Duration   time.Duration
}

// request command.
func request(url string) tea.Cmd {
	return func(ctx context.Context) tea.Msg {
		start := time.Now()

		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		return requestCompleted{
			StatusCode: res.StatusCode,
			Duration:   time.Since(start),
		}
	}
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)
	switch msg := msg.(type) {
	case requestCompleted:
		m.Iteration++
		m.Progress.Percent = float64(m.Iteration) / float64(m.MaxIterations)
		if m.Iteration >= m.MaxIterations {
			return m, tea.Quit
		}
		m.Previous = msg
		m.Durations = append(m.Durations, msg.Duration)
		return m, request(m.URL)
	case *terminput.KeyboardInput:
		// pressed esc or q
		if msg.Key() == terminput.KeyEscape || msg.Rune() == 'q' {
			return m, tea.Quit
		}
		return m, nil
	}

	return m, nil
}

// view function.
func view(ctx context.Context, model tea.Model) string {
	m := model.(Model)
	var w strings.Builder
	fmt.Fprintf(&w, "\r\n")
	if m.Progress.Percent < 1 {
		fmt.Fprintf(&w, "  Benchmarking \033[1m%s\033[m\r\n\r\n", m.URL)
		fmt.Fprintf(&w, "  %s\r\n\r\n", progress.View(m.Progress))
		fmt.Fprintf(&w, "   Request: %d of %d\r\n", m.Iteration, m.MaxIterations)
		fmt.Fprintf(&w, "    Status: %d\r\n", m.Previous.StatusCode)
		fmt.Fprintf(&w, "  Duration: %s\r\n", m.Previous.Duration.Round(time.Millisecond))
	} else {
		fmt.Fprintf(&w, "  Min: %s\r\n", min(m.Durations).Round(time.Millisecond))
		fmt.Fprintf(&w, "  Avg: %s\r\n", avg(m.Durations).Round(time.Millisecond))
		fmt.Fprintf(&w, "  Max: %s\r\n", max(m.Durations).Round(time.Millisecond))
	}
	fmt.Fprintf(&w, "\r\n")
	return w.String()
}

// min of durations.
func min(durations []time.Duration) (v time.Duration) {
	v = durations[0]
	for _, d := range durations[1:] {
		if d < v {
			v = d
		}
	}
	return
}

// max of durations.
func max(durations []time.Duration) (v time.Duration) {
	v = durations[0]
	for _, d := range durations[1:] {
		if d > v {
			v = d
		}
	}
	return
}

// sum of durations.
func sum(durations []time.Duration) (v time.Duration) {
	for _, d := range durations {
		v += d
	}
	return
}

// avg of durations.
func avg(durations []time.Duration) time.Duration {
	return sum(durations) / time.Duration(len(durations))
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\r\n", err)
	}
}
