package main

import "github.com/charmbracelet/lipgloss"

var (
	// Muted violet headline
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a387c4")). // soft lavender
			Bold(true).
			Underline(true)

	// Soft blue subheadings
	subHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7aa2f7")). // muted periwinkle
			Bold(true)

	// Gentle green for success
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ece6a")) // desaturated lime

	// Subtle gold for warnings
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0af68")). // sandy gold
			Bold(true)

	// Muted red for errors
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f7768e")). // rosy red
			Bold(true)

	// Cool teal for info messages
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7dcfff")) // soft cyan

	// Gray prefix symbol
	prefixStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")). // pastel gray-blue
			SetString("→")
)

// Maps container status to soft-styled output
func statusColor(status string) lipgloss.Style {
	switch status {
	case "running":
		return successStyle
	case "stopped":
		return warningStyle
	default:
		return errorStyle
	}
}
