package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/tj/go-tea"
	"github.com/tj/go-tea/viewport"
	"github.com/tj/go-terminput"
)

// Model struct.
type Model struct {
	List viewport.Model
}

// initialize function.
func initialize(ctx context.Context) (tea.Model, tea.Cmd) {
	return Model{
		List: viewport.Model{
			Height:       45,
			ScrollBy:     5,
			ScrollHeight: 100,
		},
	}, nil
}

// update function.
func update(ctx context.Context, msg tea.Msg, model tea.Model) (tea.Model, tea.Cmd) {
	m := model.(Model)
	m.List = viewport.Update(msg, m.List)

	switch msg := msg.(type) {
	case *terminput.KeyboardInput:
		switch msg.Key() {
		case terminput.KeyEscape:
			return m, tea.Quit
		case terminput.KeyRune:
			switch r := msg.Rune(); r {
			case 'q':
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// view function.
func view(ctx context.Context, model tea.Model) string {
	w := new(bytes.Buffer)
	m := model.(Model)

	// padding
	fmt.Fprintf(w, "\n")
	defer fmt.Fprintf(w, "\n")

	// list
	fmt.Fprintf(w, viewport.View(m.List, viewList(100)))

	return w.String()
}

// viewList returns a generated list of n items.
func viewList(n int) (s string) {
	for i := 0; i < n; i++ {
		s += fmt.Sprintf("  %d) Some blog post\n", i)
	}
	return
}

func main() {
	program := tea.NewProgram(initialize, update, view)
	err := program.Start(context.Background())
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
