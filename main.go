package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}

func main() {
	Run()
}
