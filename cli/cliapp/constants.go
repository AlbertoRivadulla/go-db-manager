package cliapp

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorRed    = lipgloss.Color("#f54242")
	colorYellow = lipgloss.Color("#b0ad09")
	colorBlue   = lipgloss.Color("#347aeb")
	colorGray   = lipgloss.Color("#636363")
	colorGreen  = lipgloss.Color("#1fb009")
	colorWhite  = lipgloss.Color("#FFFDF5")
)

var (
	whiteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRed)
	yellowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)
	grayStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGray)
	goodStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)
	blueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBlue)

	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			// Faint(true).
			Foreground(colorGray).
			Padding(0, 1)

	quitKey = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)

	backKey = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	)

	selectKey = key.NewBinding(
		key.WithKeys("enter", "o"),
		key.WithHelp("enter/o", "select"),
	)
)
