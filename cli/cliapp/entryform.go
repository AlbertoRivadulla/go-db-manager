package cliapp

import (
	"fmt"
	"strings"

	"carmaintenance/internal/core/spec-models"

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
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch {
	case key.Matches(msg, nextEntryKey):
		f.Fields[f.FocusedField].Input.Blur()
		f.FocusedField = (f.FocusedField + 1) % len(f.Fields)
		cmd = f.Fields[f.FocusedField].Input.Focus()
		return f, cmd
	case key.Matches(msg, prevEntryKey):
		f.Fields[f.FocusedField].Input.Blur()
		if f.FocusedField == 0 {
			f.FocusedField = len(f.Fields) - 1
		} else {
			f.FocusedField--
		}
		cmd = f.Fields[f.FocusedField].Input.Focus()
		return f, cmd
	}

	// Update the text in the focused field
	f.Fields[f.FocusedField].Input, cmd = f.Fields[f.FocusedField].Input.Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
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

	helpText := helpStyle.Render("Tab/↓: next field • Shift+Tab/↑: previous field • enter: add entry • esc: go back")

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

func (f EntryForm) GetValues() ([]string, []string, error) {
	var keys []string
	var values []string

	if err := f.validate(); err != nil {
		return keys, values, fmt.Errorf("could not get values from form: %w", err)
	}

	for _, field := range f.Fields {
		value := field.Input.Value()
		if value != "" {
			keys = append(keys, field.Label)
			values = append(values, value)
		}
	}

	return keys, values, nil
}

func (f EntryForm) validate() error {

	for _, field := range f.Fields {
		if err := validateField(field); err != nil {
			return fmt.Errorf("invalid field \"%s\": %w", field.Label, err)
		}
	}

	return nil
}

func validateField(field FormField) error {
	value := field.Input.Value()

	if field.Required && value == "" {
		return core.ValidationError{fmt.Sprintf("field must not be empty")}
	}

	if value != "" {
		switch field.DataType {
		case "integer":
			return core.ValidateInt(value)
		case "real":
			return core.ValidateReal(value)
		case "text":
			return core.ValidateText(value)
		case "timestamp":
			return core.ValidateTimestamp(value)
		}
	}

	return nil
}
