package cliapp

import (
    "github.com/charmbracelet/lipgloss"
)

type StatusBarState string

const (
	StatusBarStateOk StatusBarState = "green"
	StatusBarStateErr StatusBarState = "red"
)

var styleMapByColor = map[StatusBarState]lipgloss.Style{
	StatusBarStateOk:  statusStyleGreen,
	StatusBarStateErr:    statusStyleErr,
}

type StatusBarProps struct {
	State StatusBarState
	Message string
	Width int
	Count int
}

var (
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyleGreen = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color(colorGreen)).
				Padding(0, 1).
				MarginRight(1)

	statusStyleErr = statusStyleGreen.Background(colorRed)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)
)

func (s *StatusBarProps) Render() string {
	coloredStyle, ok := styleMapByColor[s.State]
	if !ok {
		coloredStyle = statusStyleGreen
	}

	statusKey := coloredStyle.Render("STATUS")

	w := lipgloss.Width
	statusVal := statusText.
		Width(s.Width - w(statusKey)).
		Render(whiteStyle.Render(s.Message))

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
	)

	return bar
}
