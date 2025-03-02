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
	orca "github.com/predixus/orca/internal"
	dlyr "github.com/predixus/orca/internal/datalayers"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
	dlyr          textinput.Model
	connStr       textinput.Model
	help          help.Model
	keys          keyMap
	err           error
	datalayerType dlyr.Platform
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
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// initialise the model
func initialModel() model {
	tiDlyr := textinput.New()
	tiDlyr.Placeholder = "PostgreSQL"
	tiDlyr.Focus()
	tiDlyr.CharLimit = 156
	tiDlyr.Width = 50

	tiConnStr := textinput.New()
	tiConnStr.Placeholder = "postgresql://orca:orca_password@orca_postgres:5432/orca"
	tiConnStr.Focus()
	tiConnStr.CharLimit = 156
	tiConnStr.Width = 50

	return model{
		state:   configuring,
		dlyr:    tiDlyr,
		connStr: tiConnStr,
		help:    help.New(),
		keys:    keys,
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
		case key.Matches(msg, m.keys.Start):
			if m.state == configuring && m.connStr.Value() != "" {
				m.state = running
				return m, startGRPCServer(m.connStr.Value())
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
		m.dlyr, cmd = m.dlyr.Update(msg)
		m.connStr, cmd = m.connStr.Update(msg)
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
		s.WriteString("\nSelect datalayer\n")
		s.WriteString(m.dlyr.View())
		// s.WriteString("\nEnter database connection string:\n")
		// s.WriteString(m.connStr.View())
		// s.WriteString("\n\n")
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
