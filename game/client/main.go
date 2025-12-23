package client

import tea "github.com/charmbracelet/bubbletea"

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

// Update implements [tea.Model].
func (m model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	panic("unimplemented")
}

// View implements [tea.Model].
func (m model) View() string {
	panic("unimplemented")
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func main() {
	tui := tea.NewProgram(initialModel())
	if _, err := tui.Run(); err != nil {
		panic(err)
	}
}
