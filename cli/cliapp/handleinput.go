package cliapp

import (
	"fmt"

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
	// m.mainMenuList, cmd = m.mainMenuList.Update(msg)
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

	switch msg.String() {
	case "enter":
		item, ok := m.mainMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}
		switch item.title {
		case "Manage tables":
			m.currScreenFocus = ViewTablesScreen
			m.currScreenLeft = ViewTablesScreen
			m.viewTablesMenuList.SetItems(m.getTablesAsMenuItems())

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

func (m Model) handleViewTablesMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// // Don't intercept keys when filtering is active
	// if m.viewTablesMenuList.FilterState() == list.Filtering {
	// 	var cmd tea.Cmd
	// 	m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	// 	return m, cmd
	// }

	switch msg.String() {
	case "enter":
		item, ok := m.viewTablesMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}

		m.currScreenFocus = ManageTableScreen
		m.currScreenLeft = ManageTableScreen
		// TODO: Do I need to set m.currScreenRight?
		m.manageTableMenuList.Title = fmt.Sprintf("Manage table - %s", item.title)
		m.currTableName = item.title
		m.currTableDesc = item.desc
		m.showTable = true

		// Set the style of the table to plain
		s := table.DefaultStyles()
		s.Selected = lipgloss.NewStyle()
		s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
		m.table.SetStyles(s)

		// Get the columns in the table
		columns, err := m.backend.GetColumnsInTable(item.title)
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error reading columns: %s", err)
			return m, nil
		}

		// Load the data from the table
		result, err := m.backend.GetEntriesInTable(item.title)
		if err != nil {
			m.status.State = StatusBarStateErr
			m.status.Message = fmt.Sprintf("Error reading table entries: %s", err)
			return m, nil
		}
		m.status.State = StatusBarStateOk
		m.status.Message = fmt.Sprintf("Read table %v", item.title)

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
	}

	var cmd tea.Cmd
	m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	return m, cmd
}

func (m Model) handleManageTableMenuScreen(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		item, ok := m.manageTableMenuList.SelectedItem().(menuItem)
		if !ok {
			return m, nil
		}

		switch item.title {
		case "View table":
			m.currScreenFocus = TableScreen
			// Set the style of the selected row when navigating it
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
			m.table.SetStyles(s)

		case "Add entry":
			var cmd tea.Cmd
			return m, cmd
			// TODO: Implement this

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
	switch msg.String() {
	case "enter":
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
