package main

import tea "github.com/charmbracelet/bubbletea"

type screen int

const (
	mainMenuScreen screen = iota
	createGameScreen
	joinGameScreen
)

type model struct {
	screen  screen   // current screen
	choices []string // menu options
	cursor  int      // which menu item our cursor is pointing at
}

// Update implements [tea.Model].
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case mainMenuScreen:
		return m.updateMainMenu(msg)
	case createGameScreen:
		return m.updateCreateGame(msg)
	case joinGameScreen:
		return m.updateJoinGame(msg)
	}
	return m, nil
}

func (m model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0: // Create Game
				m.screen = createGameScreen
				m.cursor = 0
				return m, nil
			case 1: // Join Game
				m.screen = joinGameScreen
				m.cursor = 0
				return m, nil
			case 2: // Exit
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) updateCreateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Go back to main menu
			m.screen = mainMenuScreen
			m.cursor = 0
			m.choices = []string{"Create Game", "Join Game", "Exit"}
			return m, nil
		}
	}
	return m, nil
}

func (m model) updateJoinGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Go back to main menu
			m.screen = mainMenuScreen
			m.cursor = 0
			m.choices = []string{"Create Game", "Join Game", "Exit"}
			return m, nil
		}
	}
	return m, nil
}

// View implements [tea.Model].
func (m model) View() string {
	switch m.screen {
	case mainMenuScreen:
		return m.viewMainMenu()
	case createGameScreen:
		return m.viewCreateGame()
	case joinGameScreen:
		return m.viewJoinGame()
	}
	return ""
}

func (m model) viewMainMenu() string {
	s := "╔══════════════════════════════════╗\n"
	s += "║     Welcome to StackItManGO!     ║\n"
	s += "╚══════════════════════════════════╝\n\n"

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "► "
		}
		s += cursor + choice + "\n"
	}

	s += "\n(↑/↓, enter to select)\n"
	return s
}

func (m model) viewCreateGame() string {
	s := "╔══════════════════════════════════╗\n"
	s += "║         Create New Game          ║\n"
	s += "╚══════════════════════════════════╝\n\n"
	s += "Game creation screen goes here...\n\n"
	s += "(esc to go back)\n"
	return s
}

func (m model) viewJoinGame() string {
	s := "╔══════════════════════════════════╗\n"
	s += "║           Join Game              ║\n"
	s += "╚══════════════════════════════════╝\n\n"
	s += "Join game screen goes here...\n\n"
	s += "(esc to go back)\n"
	return s
}

func initialModel() model {
	return model{
		screen:  mainMenuScreen,
		choices: []string{"Create Game", "Join Game", "Exit"},
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
