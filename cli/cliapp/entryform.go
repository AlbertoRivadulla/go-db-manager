package cliapp

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormField struct {
	Label        string
	DataType     string
	DefaultValue string
	Required     bool
	Input        textinput.Model
}

type EntryForm struct {
	Fields       []FormField
	FocusedField int

	Title      string
	TitleStyle lipgloss.Style

	Height int
}

func NewEntryForm() EntryForm {
	return EntryForm{
		FocusedField: 0,
		TitleStyle:   lipgloss.NewStyle(),
	}
}

func (f EntryForm) Update(msg tea.KeyMsg) (EntryForm, tea.Cmd) {
	switch {
	case key.Matches(msg, nextEntryKey):
		f.FocusedField = (f.FocusedField + 1) % len(f.Fields)
	case key.Matches(msg, prevEntryKey):
		if f.FocusedField == 0 {
			f.FocusedField = len(f.Fields) - 1
		} else {
			f.FocusedField--
		}

		// TODO: Key enter
	}

	// Update the text in the focused field
	var cmd tea.Cmd
	f.Fields[f.FocusedField].Input, cmd = f.Fields[f.FocusedField].Input.Update(msg)

	return f, cmd
}

func (f EntryForm) View() string {
	var sections []string

	topHeight := 0

	sections = append(sections, f.TitleStyle.Render(f.Title))
	sections = append(sections, "")
	topHeight += 2

	for i, field := range f.Fields {
		fieldString := f.renderField(field, i == f.FocusedField)
		sections = append(sections, fieldString)
		topHeight += lipgloss.Height(fieldString)
	}
	topHeight += len(sections) + 4

	helpText := helpStyle.Render("Tab/Shift+Tab: navigate • enter: accept")

	spacerHeight := f.Height - topHeight - lipgloss.Height(helpText)
	if spacerHeight > 0 {
		sections = append(sections, strings.Repeat("\n", spacerHeight))
	}

	sections = append(sections, helpText)

	content := strings.Join(sections, "\n\n")

	return content
}

func (f EntryForm) renderField(field FormField, focused bool) string {
	var b strings.Builder

	label := field.Label
	if field.Required {
		label += " *"
	}
	b.WriteString(formLabelStyle.Render(label) + "\t" + grayStyle.Render(field.DataType))
	b.WriteString("\n")

	// Input box
	var borderColor lipgloss.Color
	if focused {
		borderColor = textFocusedBorderColor
	} else {
		borderColor = textUnfocusedBorderColor
	}

	inputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(60)

	b.WriteString(inputStyle.Render(field.Input.View()))

	return b.String()
}

// TODO: Implement functions:
// 	- GetValues -> Get all the values as a map
// 	- Validate ->
