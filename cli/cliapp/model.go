package cliapp

import (
	"context"
	// "fmt"

	"carmaintenance/internal/core"

	tea "github.com/charmbracelet/bubbletea"
)

type Screen int

const (
	MainMenuScreen Screen = iota
	SelectTablesScreen
	ViewTableScreen
	AddEntryScreen
	EditEntryScreen
	SelectQueryScreen
)

type Model struct {
	screen Screen

	cursor int
	choices []string

	log []string

	ctx context.Context
	backend *core.Backend
}

func NewModel(ctx context.Context, backend *core.Backend) Model {
	return Model{
		screen: MainMenuScreen,
		ctx: ctx,
		backend: backend,
	}
}

func (m Model) Init() tea.Cmd {
	// TODO: Implement
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			// example DB interaction
			// result := queryDB(m.ctx, m.db)
			// m.log = append(m.log, result)
			m.log = append(m.log, "pressed enter")
		}

	case tea.WindowSizeMsg:
		// handle resize if needed
	}

	return m, nil
}

func (m Model) View() string {
	// TODO: Implement
	s := "Press Enter to run query\nPress q or Ctrl+C to quit\n\n"
	for _, line := range m.log {
		s += line + "\n"
	}
	return s
}
