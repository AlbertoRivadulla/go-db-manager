package cliapp

import (
	"context"

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
	TableScreen
)

type Model struct {
	currScreenFocus Screen
	currScreenLeft Screen
	currScreenRight Screen // TODO: I think I don't need this

	// Elements of the left column
	mainMenuList list.Model
	viewTablesMenuList list.Model
	manageTableMenuList list.Model
	runQueriesMenuList list.Model
	// TODO: Menu for applying the rules

	// Elements of the right column
	table table.Model
	viewport viewport.Model // TODO: I think I don't need this

	currTableName string
	currTableDesc string
	showTable bool
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
	mainMenuList.Styles.Title = titleStyle
	mainMenuList.Title = "Main menu"
	mainMenuList.InfiniteScrolling = true
	mainMenuList.SetShowStatusBar(false)
	mainMenuList.SetFilteringEnabled(false)

	viewTablesMenuList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	viewTablesMenuList.Styles.Title = titleStyle
	viewTablesMenuList.Title = "View tables"
	viewTablesMenuList.InfiniteScrolling = true
	viewTablesMenuList.SetShowStatusBar(false)
	viewTablesMenuList.SetFilteringEnabled(true)

	manageTableMenuItems := []list.Item{
		menuItem{title: "View table", desc: "Show all the entries of the table"},
		menuItem{title: "Add entry", desc: "Add a new entry to the table"},
		menuItem{title: "Edit entry", desc: "Select an entry and edit or remove it"},
	}
	manageTableMenuList := list.New(manageTableMenuItems, list.NewDefaultDelegate(), 0, 0)
	manageTableMenuList.Styles.Title = titleStyle
	manageTableMenuList.Title = "Manage table - ???? (Table title)"
	manageTableMenuList.InfiniteScrolling = true
	manageTableMenuList.SetShowStatusBar(false)
	manageTableMenuList.SetFilteringEnabled(false)
	// TODO: Set the title after selecting a table

	runQueriesMenuList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	runQueriesMenuList.Styles.Title = titleStyle
	runQueriesMenuList.Title = "Run query"
	runQueriesMenuList.InfiniteScrolling = true
	runQueriesMenuList.SetShowStatusBar(false)
	runQueriesMenuList.SetFilteringEnabled(false)

	dataTable := table.New(
		table.WithFocused(true),
		table.WithHeight(10), // TODO: Set this
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	dataTable.SetStyles(s)

	// Viewport for text results
	viewport := viewport.New(0, 0)

	status := StatusBarProps{
		State: "green",
		Message: "",
	}

	return Model{
		currScreenFocus: MainMenuScreen,
		currScreenLeft: MainMenuScreen,
		currScreenRight: ViewTablesScreen,

		mainMenuList: mainMenuList,
		viewTablesMenuList: viewTablesMenuList,
		manageTableMenuList: manageTableMenuList,
		runQueriesMenuList: runQueriesMenuList,

		viewport: viewport,
		table: dataTable,

		status: status,
		showTable: false,

		ctx: ctx,
		backend: backend,
	}
}

func (m Model) Init() tea.Cmd {
	// TODO: Implement
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if !m.filterState() {
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if !m.filterState() {
				return m, nil
			}
			// TODO: Implement history and go back to the previous screen
			// TODO: If the history is empty, close the app
		}

		return m.handleScreenInput(msg)

	case tea.WindowSizeMsg:
		// handle resize if needed
		m.width = msg.Width
		m.height = msg.Height
		m.mainItemsHeight = m.height - 2

		// Column width
		m.leftWidth = int(float64(m.width) * 0.4)
		m.rightWidth = m.width - m.leftWidth - 2

		m.mainMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.viewTablesMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.manageTableMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.runQueriesMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.viewport.Width = m.rightWidth
		m.viewport.Height = m.mainItemsHeight

		m.status.Width = m.width

		m.table.SetHeight(m.height - 6)

		if !m.ready {
			m.ready = true
		}
	}

	if m.filterState() {
		// m, cmd = m.updateFilteredScreen(msg)
		return m.updateFilteredScreen(msg)
	}

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Intializing the CLI app..."
	}
	leftColumn := m.RenderLeftColumn()
	rightColumn := m.RenderRightColumn()

	// Combine the two columns
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		rightColumn,
	)

	status := m.status.Render()

	return content + "\n" + status
}

func (m Model) filterState() bool {
	// TODO: Take into account the different screens that can be in the filtering state
	return m.viewTablesMenuList.FilterState() == list.Filtering
}

func (m Model) updateFilteredScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Take into account the different screens that can be in the filtering state
	var cmd tea.Cmd
	switch m.currScreenFocus {
	case ViewTablesScreen:
		m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	}

	return m, cmd
}
