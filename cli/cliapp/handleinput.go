package cliapp

import (
	"github.com/charmbracelet/bubbles/list"
    tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleScreenInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currScreenFocus {
	case MainMenuScreen:
		return m.handleMainMenuScreen(msg)
	case ViewTablesScreen:
		return m.handleViewTablesMenuScreen(msg)
	case ManageTableScreen:
		// TODO:
	case AddEntryScreen:
		// TODO:
	case EditEntryScreen:
		// TODO:
	case SelectQueryScreen:
		// TODO:
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
			m.viewTablesMenuList.ResetFilter()

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

func (m Model) getTablesAsMenuItems() ([]list.Item) {
	tableSpecs := m.backend.GetTableSpecs()

	menuItems := make([]list.Item, len(tableSpecs))
	for i, spec := range tableSpecs {
		menuItems[i] = menuItem{
			title: spec.Table.Name,
			desc: spec.Table.Description,
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

	// TODO:

	var cmd tea.Cmd
	m.viewTablesMenuList, cmd = m.viewTablesMenuList.Update(msg)
	return m, cmd
}
