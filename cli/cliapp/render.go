package cliapp

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
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
		return columnStyle.Render(m.addEntryForm.View())

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

		sections = append(sections, titleStyle.Width(m.rightWidth-4).Render(m.currTableName))
		sections = append(sections, "")

		sections = append(sections, descStyle.Render(m.currTableDesc))
		sections = append(sections, "")

		sections = append(sections, descStyle.Render(fmt.Sprintf("%d entries", m.currTableNrEntries)))
		sections = append(sections, "")

		sections = append(sections, m.table.View())

		// TODO: Finish this
		helpText := helpStyle.Render("↑/↓: navigate • d: delete • enter: view entry • ???")
		sections = append(sections, helpText)

		content := strings.Join(sections, "\n")

		containerStyle := lipgloss.NewStyle().
			// Border(lipgloss.RoundedBorder()).
			// BorderForeground(lipgloss.Color("62")).
			Padding(1).
			Width(m.rightWidth)

		if m.showTableModal {
			content = m.overlayTableModal(content)
		}

		return containerStyle.Render(content)
	}

	return ""
}

func (m Model) overlayTableModal(baseView string) string {
	if m.tableModalIndex < 0 || m.tableModalIndex > len(m.table.Rows()) {
		return baseView
	}

	columns := m.table.Columns()
	item := m.table.Rows()[m.tableModalIndex]

	modalContent, modalTextWidth := m.renderTableModalContent(columns, item)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(modalBorderColor).
		Background(modalBackgroundColor).
		Padding(1, 2).
		Width(modalTextWidth + 6)

	modal := modalStyle.Render(modalContent)

	output := overlay.Composite(modal, baseView, overlay.Center, overlay.Center, 0, 0)

	return output
}

func (m Model) renderTableModalContent(columns []table.Column, item table.Row) (string, int) {
	var b strings.Builder

	textWidth := 0

	for i, col := range columns {
		b.WriteString(whiteStyle.Render(col.Title))
		if len(col.Title) > textWidth {
			textWidth = len(col.Title)
		}
		b.WriteString("\n")

		if len(item[i]) > m.maxTableModalWidth {
			text, thisWidth := wrapText(item[i], m.maxTableModalWidth)
			b.WriteString(text)
			if thisWidth > textWidth {
				textWidth = thisWidth
			}

		} else {
			b.WriteString(item[i])
		}

		if i != len(columns)-1 {
			b.WriteString("\n\n")
		}
	}

	return b.String(), textWidth
}

func wrapText(text string, maxWidth int) (string, int) {
	if len(text) <= maxWidth {
		return text, len(text)
	}

	width := 0

	var lines []string
	words := strings.Fields(text)
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= maxWidth {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)

			if len(currentLine) > width {
				width = len(currentLine)
			}

			currentLine = word
		}
	}

	lines = append(lines, currentLine)
	return strings.Join(lines, "\n"), width
}
