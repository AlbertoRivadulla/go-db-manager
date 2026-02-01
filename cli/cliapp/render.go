package cliapp

import (
	"strings"

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
		return columnStyle.Render(m.manageTableMenuList.View())
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
	if m.showTable {
		var sections []string

		sections = append(sections, titleStyle.Width(m.rightWidth - 4).Render(m.currTableName))
		sections = append(sections, "")

		sections = append(sections, descStyle.Render(m.currTableDesc))
		sections = append(sections, "")

		sections = append(sections, m.table.View())

		content := strings.Join(sections, "\n")

		containerStyle := lipgloss.NewStyle().
			// Border(lipgloss.RoundedBorder()).
			// BorderForeground(lipgloss.Color("62")).
			Padding(1).
			Width(m.rightWidth)

		// columnStyle := lipgloss.NewStyle().
		//     Padding(1).Width(m.leftWidth)

		return containerStyle.Render(content)
	}

	return ""
}
