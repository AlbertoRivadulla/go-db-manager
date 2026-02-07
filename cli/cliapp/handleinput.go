package cliapp

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) handleScreenInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currScreenFocus {
	case MainMenuScreen:
		return m.handleMainMenuScreen(msg)
	case ViewTablesScreen:
		return m.handleViewTablesMenuScreen(msg)
	case ManageTableScreen:
		return m.handleManageTableMenuScreen(msg)
	case AddEntryScreen:
		// TODO:
	case EditEntryScreen:
		// TODO:
	case SelectQueryScreen:
		// TODO:
	case TableScreen:
		return m.handleTableScreen(msg)
	}

	// var cmds []tea.Cmd
	// // // TODO: Update the currently active item in the left
	// // // Modify this
	// m.m1ainMenuList, cmd = m.mainMenuList.Update(msg)
	// cmds = append(cmds, cmd)
	// //
	// // // TODO: Update the currently active item in the right
	// // // Modify this
	// // m.viewport, cmd = m.viewport.Update(msg)
	// // cmds = append(cmds, cmd)

	return m, cmd
}

func (m Model) handleMainMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// // Don't intercept keys when filtering is active
	// if m.mainMenuList.FilterState() == list.Filtering {
	// 	var cmd tea.Cmd
	// 	m.mainMenuList, cmd = m.mainMenuList.Update(msg)
	// 	return m, cmd
	// }

	switch {
	case key.Matches(msg, selectKey):
		item, ok := m.mainMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}
		switch item.title {
		case "Manage tables":
			return m.navigateTo(m.currScreenLeft, ViewTablesScreen), nil

		case "Run queries":
			var cmd tea.Cmd
			return m, cmd
			// TODO: Implement this
		}
	}

	var cmd tea.Cmd
	m.mainMenuList, cmd = m.mainMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleViewTablesMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// // Don't intercept keys when filtering is active
	// if m.viewTablesMenuList.FilterState() == list.Filtering {
	// 	var cmd tea.Cmd
	// 	m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	// 	return m, cmd
	// }

	switch {
	case key.Matches(msg, selectKey):
		item, ok := m.viewTablesMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}

		m.currTableName = item.title
		m.currTableDesc = item.desc

		return m.navigateTo(m.currScreenLeft, ManageTableScreen), nil
	}

	var cmd tea.Cmd
	m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleManageTableMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, selectKey):
		item, ok := m.manageTableMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}

		switch item.title {
		case "View table":
			return m.navigateTo(m.currScreenLeft, TableScreen), nil

		case "Add entry":
			return m.navigateTo(m.currScreenLeft, AddEntryScreen), nil

		case "Edit entry":
			var cmd tea.Cmd
			return m, cmd
			// TODO: Implement this
		}
	}

	var cmd tea.Cmd
	m.manageTableMenuList, cmd = m.manageTableMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleTableScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, selectKey):
		// item, ok := m.manageTableMenuList.SelectedItem().(menuItem)
		// if !ok {
		// 	return m, nil
		// }
		// TODO: Implement
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) NavigateBack() tea.Model {
	if len(m.screenHistory) > 0 {
		screen := m.screenHistory[len(m.screenHistory)-1]
		m.screenHistory = m.screenHistory[:len(m.screenHistory)-1]
		return m.navigateTo(None, screen)
	}

	return m
}

func (m Model) navigateTo(fromScreen, toScreen Screen) tea.Model {
	if fromScreen != None {
		m.screenHistory = append(m.screenHistory, fromScreen)
	}

	m.currScreenFocus = toScreen

	switch toScreen {
	case MainMenuScreen:
		m.currScreenLeft = toScreen
	case ViewTablesScreen:
		m.currScreenLeft = toScreen
		m.viewTablesMenuList.SetItems(m.getTablesAsMenuItems())
	case ManageTableScreen:
		m.currScreenLeft = toScreen

		err := m.setupTable()
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error setting up table %s: %s", m.currTableName, err)
			return m
		}

		// Set the style of the table to plain
		s := table.DefaultStyles()
		s.Selected = lipgloss.NewStyle()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		m.table.SetStyles(s)
	case TableScreen:
		// Set the style of the selected row when navigating it
		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = tableSelectedItemHighlight
		m.table.SetStyles(s)

	case AddEntryScreen:
		m.currScreenLeft = toScreen

		err := m.setupEntryForm()
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error setting form for table %s: %s", m.currTableName, err)
			return m
		}
	}

	return m
}

func (m Model) getTablesAsMenuItems() []list.Item {
	tableSpecs := m.backend.GetTableSpecs()

	menuItems := make([]list.Item, len(tableSpecs))
	for i, spec := range tableSpecs {
		menuItems[i] = menuItem{
			title: spec.Table.Name,
			desc:  spec.Table.Description,
		}
	}
	return menuItems
}

func (m *Model) setupTable() error {
	// TODO: Do I need to set m.currScreenRight?
	m.manageTableMenuList.Title = fmt.Sprintf("Manage table - %s", m.currTableName)
	m.showTable = true

	// Get the columns in the table
	columns, err := m.backend.GetColumnsInTable(m.currTableName)
	if err != nil {
		return fmt.Errorf("Error reading columns: %w", err)
	}

	// Load the data from the table
	result, err := m.backend.GetEntriesInTable(m.currTableName)
	if err != nil {
		return fmt.Errorf("Error reading table entries: %w", err)
	}
	m.status.State = StatusBarStateOk
	m.status.Message = fmt.Sprintf("Read table %v", m.currTableName)

	tableColumns := make([]table.Column, len(columns))
	for i, col := range columns {
		tableColumns[i] = table.Column{
			Title: col,
			Width: len(col) + 2,
		}
	}

	tableRows := make([]table.Row, len(result))

	if len(result) > 0 {
		for i, row := range result {
			tableRows[i] = row
		}
	}

	m.table.SetColumns(tableColumns)
	m.table.SetRows(tableRows)

	m.currTableNrEntries = len(tableRows)

	return nil
}

func (m *Model) setupEntryForm() error {
	columnSpecs, err := m.backend.GetColumnSpecsInTable(m.currTableName)
	if err != nil {
		return fmt.Errorf("Error getting the column specs for table %s: %w", m.currTableName, err)
	}

	m.addEntryForm.Cursor = 0
	m.addEntryForm.Title = fmt.Sprintf("Add entry to the table %s", m.currTableName)

	// Populate the elements in the form
	// m.addEntryForm.Fields = make([]FormField, len(columnSpecs))
	m.addEntryForm.Fields = []FormField{}

	for _, colSpec := range columnSpecs {
		if colSpec.PrimaryKey {
			continue
		}

		var field FormField
		field.Label = colSpec.Name
		field.DataType = colSpec.Type
		field.DefaultValue = ""
		field.Required = (!colSpec.Nullable)

		m.addEntryForm.Fields = append(m.addEntryForm.Fields, field)
	}

	return nil
}
