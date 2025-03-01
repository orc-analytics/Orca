package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	configuring state = iota
	running
)

type model struct {
	state       state
	dbInput     textinput.Model
	help        help.Model
	keys        keyMap
	err         error
	quitting    bool
	serverState string
}

type keyMap struct {
	Start key.Binding
	Quit  key.Binding
	Help  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "start server"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "postgresql://user:pass@localhost:5432/dbname"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return model{
		state:       configuring,
		dbInput:     ti,
		help:        help.New(),
		keys:        keys,
		serverState: "stopped",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Start):
			if m.state == configuring && m.dbInput.Value() != "" {
				m.state = running
				m.serverState = "running"
				return m, startGRPCServer(m.dbInput.Value())
			}
		}

	case error:
		m.err = msg
		return m, nil

	case serverStartedMsg:
		m.serverState = "running"
		return m, nil
	}

	if m.state == configuring {
		m.dbInput, cmd = m.dbInput.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder

	switch m.state {
	case configuring:
		s.WriteString("Welcome to Orca Server\n\n")
		s.WriteString("Enter database connection string:\n")
		s.WriteString(m.dbInput.View())
		s.WriteString("\n\n")
	case running:
		s.WriteString(fmt.Sprintf("Server State: %s\n", m.serverState))
		s.WriteString(fmt.Sprintf("Database: %s\n", m.dbInput.Value()))
	}

	if m.err != nil {
		s.WriteString(fmt.Sprintf("\nError: %v\n", m.err))
	}

	s.WriteString("\n")
	s.WriteString(m.help.View(m.keys))

	return s.String()
}

// Server management
type serverStartedMsg struct{}

func startGRPCServer(dbConnString string) tea.Cmd {
	return func() tea.Msg {
		// Here you would start your actual gRPC server
		// For now just return that we "started"
		return serverStartedMsg{}
	}
}
