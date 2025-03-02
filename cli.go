package main

import (
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	orca "github.com/predixus/orca/internal"
	dlyr "github.com/predixus/orca/internal/datalayers"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type connStringTemplate struct {
	prefix     string
	components []string
}

var connectionTemplates = map[string]connStringTemplate{
	"postgresql": {
		prefix:     "postgresql://",
		components: []string{"user", "password", "host:port", "dbname"},
	},
	"mysql": {
		prefix:     "mysql://",
		components: []string{"user", "password", "host:port", "dbname"},
	},
	"mongodb": {
		prefix:     "mongodb+srv://",
		components: []string{"user", "password", "cluster.mongodb.net", "dbname"},
	},
}

// Style for the placeholder text
var placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

// cli/server state management
type state int

const (
	quitting state = iota
	configuring
	running
	closing
)

type model struct {
	state         state
	currentInput  int // 0 for datalayer, 1 for connection string
	dlyr          textinput.Model
	connStr       textinput.Model
	help          help.Model
	keys          keyMap
	err           error
	datalayerType dlyr.Platform
}

var datalayerSuggestions = []string{
	"postgresql",
	"mysql",
	"mongodb",
}

var connectionPrefixes = map[string]string{
	"postgresql": "postgresql://",
	"mysql":      "mysql://",
	"mongodb":    "mongo+srv://",
}

// Custom TUI messages
type serverStartedMsg struct{}

// Custom key bindings
type keyMap struct {
	Start        key.Binding
	Quit         key.Binding
	Help         key.Binding
	Autocomplete key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.Help, k.Quit},
	}
}

// initialise the key map
var keys = keyMap{
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "submit answer"),
	),
	Autocomplete: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("⇥", "autocomplete"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "ctrl+q"),
		key.WithHelp("⌘+q", "quit"),
	),
}

// initialise the model
func initialModel() model {
	tiDlyr := textinput.New()
	tiDlyr.Placeholder = "postgresql"
	tiDlyr.Focus()
	tiDlyr.CharLimit = 156
	tiDlyr.Width = 50
	tiDlyr.ShowSuggestions = true
	tiDlyr.SetSuggestions(datalayerSuggestions)

	tiConnStr := textinput.New()
	tiConnStr.Placeholder = "user:password@localhost:5432/dbname"
	tiConnStr.CharLimit = 156
	tiConnStr.Width = 50

	return model{
		state:        configuring,
		currentInput: 0,
		dlyr:         tiDlyr,
		connStr:      tiConnStr,
		help:         help.New(),
		keys:         keys,
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
			m.state = quitting
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case msg.Type == tea.KeyEnter:
			if m.state == configuring {
				if m.currentInput == 0 {
					// Move from datalayer to connection string
					m.currentInput = 1
					m.dlyr.Blur()
					m.connStr.Focus()
					// Set connection string prefix based on datalayer
					if template, ok := connectionTemplates[m.dlyr.Value()]; ok {
						m.connStr.SetValue(template.prefix)
					}
				} else if m.connStr.Value() != "" {
					// Start the server
					m.state = running
					return m, startGRPCServer(m.connStr.Value())
				}
			}
		case msg.Type == tea.KeyEsc:
			if m.currentInput == 1 {
				// Go back to datalayer selection
				m.currentInput = 0
				m.connStr.Blur()
				m.dlyr.Focus()
			}
		}

	case error:
		m.err = msg
		return m, nil

	case serverStartedMsg:
		m.state = running
		return m, nil
	}

	if m.state == configuring {
		if m.currentInput == 0 {
			m.dlyr, cmd = m.dlyr.Update(msg)
		} else {
			m.connStr, cmd = m.connStr.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder
	switch m.state {
	case configuring:
		s.WriteString("------------------------- ORCA ------------------------\n")
		s.WriteString("The Orchestrated Robust-Compute and Analytics Framework\n")
		s.WriteString("-------------------------------------------------------\n")

		s.WriteString("\nSelect datalayer: \n")
		s.WriteString(m.dlyr.View())

		if m.currentInput == 1 {
			s.WriteString("\n\nEnter connection string (ESC to go back):\n")

			// Get the template for the selected datalayer
			if template, ok := connectionTemplates[m.dlyr.Value()]; ok {
				currentValue := m.connStr.Value()
				// cursor := strings.Repeat(" ", m.connStr.) + "█"
				cursor := ""

				if currentValue == template.prefix {
					// Show full template at start
					s.WriteString(template.prefix)
					s.WriteString(cursor)
					for i, comp := range template.components {
						if i > 0 {
							if i == 2 {
								s.WriteString("@")
							} else if i == 3 {
								s.WriteString("/")
							} else {
								s.WriteString(":")
							}
						}
						s.WriteString(placeholderStyle.Render("<" + comp + ">"))
					}
					s.WriteString("\n")
				} else {
					// Show current value + remaining template
					s.WriteString(currentValue)
					s.WriteString(cursor)

					// Calculate remaining parts based on separators
					input := strings.TrimPrefix(currentValue, template.prefix)
					parts := strings.FieldsFunc(input, func(r rune) bool {
						return r == ':' || r == '@' || r == '/'
					})

					if len(parts) < len(template.components) {
						remaining := template.components[len(parts):]
						nextSep := ":"
						if len(parts) == 1 {
							nextSep = "@"
						} else if len(parts) == 2 {
							nextSep = "/"
						}
						s.WriteString(nextSep)
						s.WriteString(placeholderStyle.Render("<" + remaining[0] + ">"))

						for i, comp := range remaining[1:] {
							if i == 0 && len(parts) == 1 {
								s.WriteString("/")
							} else {
								s.WriteString(":")
							}
							s.WriteString(placeholderStyle.Render("<" + comp + ">"))
						}
					}
					s.WriteString("\n")
				}
			} else {
				s.WriteString(m.connStr.View())
			}
		}
		s.WriteString("\n")
	case running:
		s.WriteString(fmt.Sprintf("Server State: %s\n", m.state))
		s.WriteString(fmt.Sprintf("Database: %s\n", m.connStr.Value()))
	}

	if m.err != nil {
		s.WriteString(fmt.Sprintf("\nError: %v\n", m.err))
	}

	s.WriteString("\n")
	s.WriteString(m.help.View(m.keys))

	return s.String()
}

func startGRPCServer(dbConnString string) tea.Cmd {
	port := 4040
	slog.Debug("Running the server", "port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		slog.Error("failed to listen", "message", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterOrcaCoreServer(grpcServer, orca.NewServer())
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		slog.Error("failed to serve", "error", err)
	}
	return func() tea.Msg {
		return serverStartedMsg{}
	}
}
