package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// valid datalayers - as they are displayed
var datalayerSuggestions = []string{
	"postgresql",
}
var currentDatalayer = "postgresql"

// templates for filling out connection string
type (
	ConnectionStrParser func(connectionStr string, example string) (map[string]string, error)
	connStringTemplate  struct {
		validationFunc ConnectionStrParser
		exampleConnStr string
	}
)

var connectionTemplates = map[string]connStringTemplate{
	"postgresql": {
		validationFunc: ParsePostgresURL,
		exampleConnStr: "postgresql://<user>:<pass>@<localhost>:<port>/<db>?<setting=value>",
	},
}

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
	port
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
	port       textinput.Model
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

// validation functions
func ValidateDatalayer(s string) error {
	if s == "" {
		return fmt.Errorf("Select a datalayer")
	}
	for _, v := range datalayerSuggestions {
		if s == v {
			currentDatalayer = v
			return nil
		}
	}
	return fmt.Errorf("Unsuported datalayer: %s", s)
}

func ValidateConnStr(s string) error {
	if s == "" {
		return errors.New("Connection string cannot be empty")
	}
	template, ok := connectionTemplates[currentDatalayer]
	if !ok { // should never occur
		return fmt.Errorf("no template found for datalayer: %s", currentDatalayer)
	}
	_, err := template.validationFunc(s, template.exampleConnStr)
	return err
}

func ValidatePort(s string) error {
	if s == "" {
		return errors.New("You have to select a port number")
	}

	// try to lookup the port to validate it
	if _, err := net.LookupPort("tcp", s); err != nil {
		return fmt.Errorf("Invalid port number '%s' (must be between 1-65535)", s)
	}

	// check if port is already in use
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s))
	if err != nil {
		return fmt.Errorf("Port %s is already in use", s)
	}
	listener.Close()

	return nil
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
	tiDlyr.Validate = ValidateDatalayer

	tiDlyr.SetSuggestions(datalayerSuggestions)

	// connection string placeholder
	tiConnStr := textinput.New()
	tiConnStr.Placeholder = connectionTemplates[defaultDlyr].exampleConnStr
	tiConnStr.CharLimit = 150
	tiConnStr.Width = 80
	tiConnStr.Validate = ValidateConnStr

	// port selection
	tiPort := textinput.New()
	tiPort.Placeholder = "4040"
	tiPort.CharLimit = 6
	tiPort.Width = 6
	tiPort.Validate = ValidatePort

	return model{
		state:      configuring,
		configStep: datalayer,
		dlyr:       tiDlyr,
		connStr:    tiConnStr,
		port:       tiPort,
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
					// move on to port selection
					m.configStep = port
					m.connStr.Blur()
					m.port.Focus()
				} else if m.configStep == port {
					if err := m.port.Validate(m.port.Value()); err != nil {
						m.err = err
						return m, nil
					} else {
						m.err = nil
					}
					// move on
					m.port.Blur()

					port, _ := strconv.Atoi(m.port.Value())
					startGRPCServer(m.connStr.Value(), port)
					m.state = running
					return m, nil
				}
			}
		}

	case error:
		m.err = msg
		return m, nil

	}

	if m.state == configuring {
		if m.configStep == datalayer {
			m.dlyr, cmd = m.dlyr.Update(msg)
			m.dlyr.SetValue(strings.ToLower(m.dlyr.Value()))
		} else if m.configStep == connectionStr {
			m.connStr, cmd = m.connStr.Update(msg)
			m.connStr.SetValue(strings.ToLower(m.connStr.Value()))
		} else if m.configStep == port {
			m.port, cmd = m.port.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	var s strings.Builder
	switch m.state {
	case configuring:
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			BorderStyle(lipgloss.RoundedBorder()).
			Align(lipgloss.Center).
			Padding(0, 1)

		subtitleStyle := lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Align(lipgloss.Center)

		bannerStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(1, 0).
			Width(60).
			Align(lipgloss.Center)

		title := titleStyle.Render("🐋 ORCA")
		subtitle := subtitleStyle.Render("The Orchestrated Robust-Compute and Analytics Framework")

		banner := bannerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Center,
				title,
				subtitle,
			),
		)

		s.WriteString(banner + "\n")

		if m.configStep == datalayer {
			s.WriteString("\nSelect datalayer: \n")
			s.WriteString(m.dlyr.View())
		}

		if m.configStep == connectionStr {
			s.WriteString("\nEnter connection string: \n")
			s.WriteString(m.connStr.View())
		}

		if m.configStep == port {
			s.WriteString("\nSelect a port number for the Orca server: \n")
			s.WriteString(m.port.View())
		}
		s.WriteString("\n")

	case running:
		style := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

		msg := style.Render(
			fmt.Sprintf("\n🐋 ORCA Server Running at grpc://localhost:%v\n", m.port.Value()),
		)
		s.WriteString(msg)
	}

	if m.err != nil {
		s.WriteString("\n" + m.err.Error() + "\n")
	}

	s.WriteString("\n")
	s.WriteString(m.help.View(m.keys))

	return s.String()
}
