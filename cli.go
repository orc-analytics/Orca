package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	orca "github.com/predixus/orca/internal"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// valid datalayers - as they are displayed
var datalayerSuggestions = []string{
	"postgresql",
}
var currentDatalayer = "postgresql"

// templates for filling out connection string
type connStringTemplate struct {
	regex          string
	exampleConnStr string
}

var connectionTemplates = map[string]connStringTemplate{
	"postgresql": {
		regex:          `(postgresql|postgres):\/\/([^:@\s]*(?::[^@\s]*)?@)?([^\/\?\s]+)`,
		exampleConnStr: "postgresql://<user>:<pass>@localhost:5432/orca?sslmode=prefer",
	},
}

// style for the placeholder text
var placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

// cli/server state management
type state int

const (
	quitting state = iota
	configuring
	running
	closing
)

// configuration steps - what are we currently configuring?
type configStep int

// this is the config step order
const (
	datalayer configStep = iota
	connectionStr
)

// the cli model
type model struct {
	state      state
	configStep configStep
	help       help.Model
	keys       keyMap
	err        error
	dlyr       textinput.Model
	connStr    textinput.Model
}

// custom TUI messages
type serverStartedMsg struct{}

// custom key bindings
type keyMap struct {
	Enter        key.Binding
	Quit         key.Binding
	Help         key.Binding
	Autocomplete key.Binding
	Esc          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Help, k.Quit, k.Esc},
	}
}

// initialise the key map
var keys = keyMap{
	Enter: key.NewBinding(
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
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("ESC", "go back"),
	),
}

// initialise the model
func initialModel() model {
	defaultDlyr := datalayerSuggestions[0]

	// datalayer selection
	tiDlyr := textinput.New()
	tiDlyr.Placeholder = defaultDlyr
	tiDlyr.Focus()
	tiDlyr.CharLimit = 150
	tiDlyr.Width = 50
	tiDlyr.ShowSuggestions = true
	tiDlyr.Validate = func(s string) error {
		if s == "" {
			return nil
		}
		for _, v := range datalayerSuggestions {
			if s == v {
				currentDatalayer = v
				return nil
			}
		}
		return fmt.Errorf("unsupported datalayer: %s", s)
	}
	tiDlyr.SetSuggestions(datalayerSuggestions)

	// connection string placeholder
	tiConnStr := textinput.New()
	tiConnStr.Placeholder = connectionTemplates[defaultDlyr].exampleConnStr
	tiConnStr.CharLimit = 150
	tiConnStr.Width = 80
	tiConnStr.Validate = func(s string) error {
		if s == "" {
			return errors.New("Datalayer string cannot be empty")
		}
		template, ok := connectionTemplates[currentDatalayer]
		if !ok { // should never enter
			return fmt.Errorf("no template found for datalayer: %s", currentDatalayer)
		}
		matched, err := regexp.Match(template.regex, []byte(s))
		if err != nil {
			return fmt.Errorf("regex error: %v", err)
		}
		if !matched {
			return fmt.Errorf("invalid connection string format")
		}
		return nil
	}

	return model{
		state:      configuring,
		configStep: datalayer,
		dlyr:       tiDlyr,
		connStr:    tiConnStr,
		help:       help.New(),
		keys:       keys,
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
		case key.Matches(msg, m.keys.Esc):
			if m.state == configuring {
				if m.configStep == connectionStr {
					m.configStep = datalayer
					m.connStr.Blur()
					m.dlyr.Focus()
				}
			}

		case key.Matches(msg, m.keys.Quit):
			m.state = quitting
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case msg.Type == tea.KeyEnter:
			if m.state == configuring {
				if m.configStep == datalayer {
					// validate datalayer selection
					if err := m.dlyr.Validate(m.dlyr.Value()); err != nil {
						m.err = err
						return m, nil
					} else {
						m.err = nil
					}

					// move from datalayer to connection string
					m.configStep = connectionStr
					m.dlyr.Blur()
					m.connStr.Focus()

				} else if m.configStep == connectionStr {
					if err := m.connStr.Validate(m.connStr.Value()); err != nil {
						m.err = err
						return m, nil
					} else {
						m.err = nil
					}

					// start the server
					m.state = running
					return m, startGRPCServer(m.connStr.Value())
				}
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
		if m.configStep == datalayer {
			m.dlyr, cmd = m.dlyr.Update(msg)
			m.dlyr.SetValue(strings.ToLower(m.dlyr.Value()))
		} else if m.configStep == connectionStr {
			m.connStr, cmd = m.connStr.Update(msg)
			m.connStr.SetValue(strings.ToLower(m.connStr.Value()))
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

		if m.configStep == connectionStr {
			s.WriteString("\n\nEnter connection string:\n")
			s.WriteString(m.connStr.View())
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
