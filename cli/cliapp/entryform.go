package cliapp

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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
	Fields []FormField
	Cursor int

	Title      string
	TitleStyle lipgloss.Style
}

func NewEntryForm() EntryForm {
	return EntryForm{
		Cursor:     0,
		TitleStyle: lipgloss.NewStyle(),
	}
}

func (f EntryForm) View() string {
	var sections []string

	sections = append(sections, f.TitleStyle.Render(f.Title))
	sections = append(sections, "")

	for i, field := range f.Fields {
		fieldStr := ""
		fieldStr += whiteStyle.Render(field.Label) + "\t" + grayStyle.Render(field.DataType) + "\n"

		// TODO: Draw the
		if i == f.Cursor {
			fieldStr += textFocusedStyle.Render(field.Input.View())
		} else {
			fieldStr += textBlurredStyle.Render(field.Input.View())
		}

		sections = append(sections, fieldStr)
	}

	content := strings.Join(sections, "\n\n")

	return content
}
