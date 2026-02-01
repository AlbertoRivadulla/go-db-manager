package cliapp

import (
	"context"
	// "strings"

	"carmaintenance/internal/core"

	"github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type Screen int

const (
	MainMenuScreen Screen = iota
	ViewTablesScreen
	ManageTableScreen
	AddEntryScreen
	EditEntryScreen
	SelectQueryScreen
)

type Model struct {
	currScreen Screen

	// Elements of the left column
	mainMenuList list.Model
	viewTablesMenuList list.Model
	manageTableMenuList list.Model
	runQueriesMenuList list.Model
	cursor int
	// TODO: Menu for applying the rules

	// Elements of the right column
	table table.Model
	viewport viewport.Model

	status StatusBarProps

	width int
	leftWidth int
	rightWidth int
	height int
	mainItemsHeight int

	ready bool

	log []string

	ctx context.Context
	backend *core.Backend
}

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

func NewModel(ctx context.Context, backend *core.Backend) Model {
	mainMenuItems := []list.Item{
		menuItem{title: "Manage tables", desc: "Operate on the available tables"},
		menuItem{title: "Run queries", desc: "Run one of the specified queries"},
	}
	mainMenuList := list.New(mainMenuItems, list.NewDefaultDelegate(), 0, 0)
	mainMenuList.Title = "Main menu"
	mainMenuList.InfiniteScrolling = true
	mainMenuList.SetShowStatusBar(false)
	mainMenuList.SetFilteringEnabled(false)

	manageTableMenuItems := []list.Item{
		menuItem{title: "View table", desc: "Show all the entries of the table"},
		menuItem{title: "Add entry", desc: "Add a new entry to the table"},
		menuItem{title: "Edit entry", desc: "Select an entry and edit or remove it"},
	}
	manageTableMenuList := list.New(manageTableMenuItems, list.NewDefaultDelegate(), 0, 0)
	manageTableMenuList.Title = "Main menu"
	manageTableMenuList.InfiniteScrolling = true
	manageTableMenuList.SetShowStatusBar(false)
	manageTableMenuList.SetFilteringEnabled(true)
	// TODO: Set the title after selecting a table

	// Viewport for text results
	viewport := viewport.New(0, 0)

	status := StatusBarProps{
		State: "green",
		Message: "",
	}

	return Model{
		currScreen: MainMenuScreen,

		mainMenuList: mainMenuList,
		manageTableMenuList: manageTableMenuList,

		viewport: viewport,

		status: status,

		ctx: ctx,
		backend: backend,
	}
}

func (m Model) Init() tea.Cmd {
	// TODO: Implement
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TODO: Check if the current screen is in filter state, and ignore the keys
		// Modify this
		if m.mainMenuList.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			// TODO: Manage interaction with the different menus
			m.status.Message = "pressed enter"
			m.status.State = "red"
		}

	case tea.WindowSizeMsg:
		// handle resize if needed
		m.width = msg.Width
		m.height = msg.Height
		m.mainItemsHeight = m.height - 2

		// Column width
		m.leftWidth = int(float64(m.width) * 0.4)
		m.rightWidth = m.width - m.leftWidth - 2

		m.mainMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.manageTableMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.viewport.Width = m.rightWidth
		m.viewport.Height = m.mainItemsHeight

		m.status.Width = m.width

		if !m.ready {
			m.ready = true
		}
	}

	// TODO: Update the currently active item in the left
	// Modify this
	m.mainMenuList, cmd = m.mainMenuList.Update(msg)
	cmds = append(cmds, cmd)

	// TODO: Update the currently active item in the right
	// Modify this
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Intializing the CLI app..."
	}

	// Define styles
    columnStyle := lipgloss.NewStyle().
        Padding(1)

    leftColumnStyle := columnStyle.Width(int(float64(m.width) * 0.4))

    rightColumnStyle := columnStyle.Width(m.width - int(float64(m.width)*0.4) - 2)

	// TODO: Render left column
	// Modify this
    leftColumn := leftColumnStyle.Render(m.mainMenuList.View())

	// TODO: Render left column
	// Modify this
    rightColumn := rightColumnStyle.Render(m.viewport.View())

	// Combine the two columns
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		rightColumn,
	)

	status := m.status.Render()

	return content + "\n" + status
}
