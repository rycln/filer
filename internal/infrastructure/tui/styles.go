package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			PaddingBottom(1)

	progressStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			PaddingBottom(1)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("156")).
			Italic(true).
			PaddingBottom(1)

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46")).
			PaddingBottom(1)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")).
			PaddingBottom(1)

	processingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			SetString("┃")
)
