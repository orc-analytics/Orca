package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

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

func init() {
	// Check for color support and set appropriate profile
	setupColorProfile()
}

// setupColorProfile detects terminal capabilities and sets appropriate color profile
func setupColorProfile() {
	// Check for explicit no-color requests
	if os.Getenv("NO_COLOR") != "" {
		lipgloss.SetColorProfile(termenv.Ascii)
		return
	}

	// Check for dumb terminal
	term := strings.ToLower(os.Getenv("TERM"))
	if term == "dumb" || term == "" {
		lipgloss.SetColorProfile(termenv.Ascii)
		return
	}

	// Check for basic terminal types that might not support colors well
	basicTerms := []string{"linux", "console", "vt100", "vt102", "vt220"}
	for _, basicTerm := range basicTerms {
		if strings.Contains(term, basicTerm) {
			lipgloss.SetColorProfile(termenv.ANSI)
			return
		}
	}

	// For CI environments, use basic colors
	if os.Getenv("CI") != "" {
		lipgloss.SetColorProfile(termenv.ANSI)
		return
	}

	// Default: let lipgloss auto-detect
	// It will choose the best profile based on terminal capabilities
}

// safeRender safely renders text with styling, falling back to plain text on error
func safeRender(style lipgloss.Style, text string) string {
	defer func() {
		if r := recover(); r != nil {
			// If styling fails, just return the plain text
			fmt.Fprintf(os.Stderr, "Warning: styling failed, using plain text\n")
		}
	}()

	// Try to render with style
	return style.Render(text)
}

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

// Helper functions for safe rendering of common styles
func renderSuccess(text string) string {
	return safeRender(successStyle, text)
}

func renderError(text string) string {
	return safeRender(errorStyle, text)
}

func renderWarning(text string) string {
	return safeRender(warningStyle, text)
}

func renderInfo(text string) string {
	return safeRender(infoStyle, text)
}

func renderHeader(text string) string {
	return safeRender(headerStyle, text)
}

func renderSubHeader(text string) string {
	return safeRender(subHeaderStyle, text)
}

// debugColorProfile prints current color profile information for debugging
func debugColorProfile() {
	fmt.Printf("TERM: %s\n", os.Getenv("TERM"))
	fmt.Printf("COLORTERM: %s\n", os.Getenv("COLORTERM"))
	fmt.Printf("NO_COLOR: %s\n", os.Getenv("NO_COLOR"))
	fmt.Printf("CI: %s\n", os.Getenv("CI"))

	profile := lipgloss.ColorProfile()
	var profileName string
	switch profile {
	case termenv.Ascii:
		profileName = "Ascii (no colors)"
	case termenv.ANSI:
		profileName = "ANSI (16 colors)"
	case termenv.ANSI256:
		profileName = "ANSI256 (256 colors)"
	case termenv.TrueColor:
		profileName = "TrueColor (24-bit)"
	default:
		profileName = "Unknown"
	}
	fmt.Printf("Color Profile: %s\n", profileName)
}
