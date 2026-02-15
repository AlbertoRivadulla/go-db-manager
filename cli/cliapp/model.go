package cliapp

import (
	"context"

	"dbmanager/internal/core"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen int

const (
	MainMenuScreen Screen = iota
	ViewTablesScreen
	ManageTableScreen
	AddEntryScreen
	SelectQueryScreen
	TableScreen
	None
)

type Model struct {
	currScreenFocus Screen
	currScreenLeft  Screen
	screenHistory   []Screen

	// Elements of the left column
	mainMenuList        list.Model
	viewTablesMenuList  list.Model
	manageTableMenuList list.Model
	runQueriesMenuList  list.Model
	addEntryForm        EntryForm

	// Elements of the right column
	table table.Model

	currTableName        string
	currTableDesc        string
	currTableNrEntries   int
	showTable            bool
	showTableModal       bool
	tableModalIndex      int
	showConfirmModal     bool
	confirmModalSelected bool
	confirmModalText     string
	confirmModalClosure  func() error
	minConfirmModalWidth int
	editEntryMode        bool
	currRowFilterCol     string
	currRowFilterVal     string
	status               StatusBarProps

	width              int
	leftWidth          int
	rightWidth         int
	maxTableModalWidth int
	height             int
	mainItemsHeight    int

	ready bool

	log []string

	ctx     context.Context
	backend *core.Backend
}

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

func NewModel(ctx context.Context, backend *core.Backend) Model {
	mainMenuItems := []list.Item{
		menuItem{title: "Manage tables", desc: "Operate on the available tables"},
		// TODO: Change this description
		menuItem{title: "Run queries", desc: "NOT IMPLEMENTED\nRun one of the specified queries"},
	}
	mainMenuList := list.New(mainMenuItems, list.NewDefaultDelegate(), 0, 0)
	mainMenuList.Styles.Title = titleStyle
	mainMenuList.Title = "Main menu"
	mainMenuList.InfiniteScrolling = true
	mainMenuList.SetShowStatusBar(false)
	mainMenuList.SetFilteringEnabled(false)
	mainMenuList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}
	mainMenuList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}

	viewTablesMenuList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	viewTablesMenuList.Styles.Title = titleStyle
	viewTablesMenuList.Title = "View tables"
	viewTablesMenuList.InfiniteScrolling = true
	viewTablesMenuList.SetShowStatusBar(true)
	viewTablesMenuList.SetFilteringEnabled(true)
	viewTablesMenuList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}
	viewTablesMenuList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}

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
	manageTableMenuList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}
	manageTableMenuList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}

	runQueriesMenuList := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	runQueriesMenuList.Styles.Title = titleStyle
	runQueriesMenuList.Title = "Run query"
	runQueriesMenuList.InfiniteScrolling = true
	runQueriesMenuList.SetShowStatusBar(false)
	runQueriesMenuList.SetFilteringEnabled(false)
	runQueriesMenuList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}
	runQueriesMenuList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{backKey, selectKey}
	}

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
	s.Selected = lipgloss.NewStyle()
	dataTable.SetStyles(s)
	// TODO:
	// runQueriesMenuList.AdditionalShortHelpKeys = func() []key.Binding {
	// 	return []key.Binding{backKey}
	// }
	// runQueriesMenuList.AdditionalFullHelpKeys = func() []key.Binding {
	// 	return []key.Binding{backKey}
	// }

	addEntryForm := NewEntryForm()
	addEntryForm.TitleStyle = titleStyle

	status := StatusBarProps{
		State:   "green",
		Message: "",
	}

	return Model{
		currScreenFocus: MainMenuScreen,
		currScreenLeft:  MainMenuScreen,

		mainMenuList:        mainMenuList,
		viewTablesMenuList:  viewTablesMenuList,
		manageTableMenuList: manageTableMenuList,
		runQueriesMenuList:  runQueriesMenuList,
		addEntryForm:        addEntryForm,

		table:            dataTable,
		showTableModal:   false,
		showConfirmModal: false,

		status:    status,
		showTable: false,

		ctx:     ctx,
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
		switch {
		case key.Matches(msg, quitKey):
			if !m.filterState() && m.currScreenFocus != AddEntryScreen {
				return m, tea.Quit
			}
		case key.Matches(msg, backKey):
			if m.showTableModal {
				m.showTableModal = false
				return m, nil
			} else if m.showConfirmModal {
				m.showConfirmModal = false
				return m, nil
			} else if m.editEntryMode {
				if m.currScreenLeft == AddEntryScreen {
					m = m.NavigateBack().(Model)
					return m.navigateTo(m.currScreenLeft, TableScreen), nil
				}
				m.editEntryMode = false
				return m.NavigateBack(), nil
			}
			if !m.filterState() {
				return m.NavigateBack(), nil
			}
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
		m.maxTableModalWidth = m.rightWidth - 40
		m.minConfirmModalWidth = 40

		m.mainMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.viewTablesMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.manageTableMenuList.SetSize(m.leftWidth, m.mainItemsHeight)
		m.runQueriesMenuList.SetSize(m.leftWidth, m.mainItemsHeight)

		m.status.Width = m.width

		m.table.SetHeight(m.height - 10)

		m.addEntryForm.Height = m.height

		if !m.ready {
			m.ready = true
		}
	}

	if m.filterState() {
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
	return m.viewTablesMenuList.FilterState() == list.Filtering
}

func (m Model) updateFilteredScreen(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currScreenFocus {
	case ViewTablesScreen:
		m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	}

	return m, cmd
}
