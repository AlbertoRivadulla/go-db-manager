package cliapp

import (
    // tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

func (m Model) RenderLeftColumn() string {
    columnStyle := lipgloss.NewStyle().
        Padding(1).Width(m.leftWidth)

	switch m.currScreenLeft {
	case MainMenuScreen:
	   return columnStyle.Render(m.mainMenuList.View())
	case ViewTablesScreen:
	   return columnStyle.Render(m.viewTablesMenuList.View())
	case ManageTableScreen:
		// TODO:
	case AddEntryScreen:
		// TODO:
	case EditEntryScreen:
		// TODO:
	case SelectQueryScreen:
		// TODO:
	}

	return ""
}

func (m Model) RenderRightColumn() string {
	//    columnStyle := lipgloss.NewStyle().
	//        Padding(1).Width(m.rightWidth)
	//
	// switch m.currScreenRight {
	// // TODO:
	// // I think I will always show the table here
	// }
	//
	// return columnStyle.Render(m.viewport.View())

	return ""
}
