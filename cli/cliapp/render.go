package cliapp

import (
	"fmt"
	"strings"

	"dbmanager/internal/utils"

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

		if m.editEntryMode {
			helpText := helpStyle.Render(fmt.Sprintf("↑/↓: navigate • %s: edit • %s: delete • %s: view entry • %s: undo",
				editEntryKey.Keys()[0],
				deleteEntryKey.Keys()[0],
				selectKey.Keys()[0],
				undoKey.Keys()[0],
			))
			sections = append(sections, helpText)
		} else {
			helpText := helpStyle.Render(fmt.Sprintf("↑/↓: navigate • %s: view entry • ???", selectKey.Keys()[0]))
			sections = append(sections, helpText)
		}

		content := strings.Join(sections, "\n")

		containerStyle := lipgloss.NewStyle().
			// Border(lipgloss.RoundedBorder()).
			// BorderForeground(lipgloss.Color("62")).
			Padding(1).
			Width(m.rightWidth)

		if m.showTableModal {
			content = m.overlayTableModal(content)
		} else if m.showConfirmModal {
			content = m.overlayConfirmModal(content)
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
			text, thisWidth := utils.WrapText(item[i], m.maxTableModalWidth)
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

func (m Model) overlayConfirmModal(baseView string) string {
	modalContent := m.renderConfirmModalContent()

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(modalBorderColor).
		Background(modalBackgroundColor).
		Padding(1, 2)

	modal := modalStyle.Render(modalContent)

	output := overlay.Composite(modal, baseView, overlay.Center, overlay.Center, 0, 0)

	return output
}

// func (m Model) renderConfirmModalContent() (string, int) {
func (m Model) renderConfirmModalContent() string {
	noButton := renderButton("No", !m.confirmModalSelected)
	yesButton := renderButton("Yes", m.confirmModalSelected)

	mainText, mainTextWidth := utils.WrapText(m.confirmModalText, m.maxTableModalWidth)
	mainText = whiteStyle.Render(mainText)
	modalWidth := max(mainTextWidth + 2, m.minConfirmModalWidth)

	emptyLineStyle := lipgloss.NewStyle().
		Background(modalBackgroundColor)
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		noButton,
		emptyLineStyle.Height(3).Render("  "),
		yesButton,
	)

	buttonsStyle := lipgloss.NewStyle().
		Background(modalBackgroundColor).
		Width(modalWidth).
		Align(lipgloss.Center)
	
	buttons = buttonsStyle.Render(buttons)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		mainText,
		emptyLineStyle.Width(modalWidth).Render(""),
		buttons,
	)
}

func renderButton(text string, selected bool) string {
	if selected {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(textFocusedBorderColor).
			BorderBackground(modalBackgroundColor).
			Background(modalBackgroundColor).
			Padding(0, 2).
			Bold(true).
			Render(text)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(textUnfocusedBorderColor).
		BorderBackground(modalBackgroundColor).
		Background(modalBackgroundColor).
		Padding(0, 2).
		Render(text)
}
