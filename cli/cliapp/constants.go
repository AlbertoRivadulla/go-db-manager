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
			Foreground(colorGray).
			Padding(0, 1)

	tableSelectedItemHighlight = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)

	// Style for the text forms
	formLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))
	textUnfocusedBorderColor = lipgloss.Color("240")
	textFocusedBorderColor = lipgloss.Color("170")
	textFocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
	textBlurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Style for modals
	modalBorderColor = lipgloss.Color("170")
	modalBackgroundColor = lipgloss.Color("235")

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

	nextEntryKey = key.NewBinding(
		key.WithKeys("tab", "down"),
		key.WithHelp("tab", "next entry"),
	)
	prevEntryKey = key.NewBinding(
		key.WithKeys("shift+tab", "up"),
		key.WithHelp("shift+tab", "previous entry"),
	)
	sendFormKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "add entry"),
	)
)
