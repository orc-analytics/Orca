package main

import (
	"context"
	"fmt"
	"log/slog"

	// "net"
	// "os"

	tea "github.com/charmbracelet/bubbletea"
	cli "github.com/predixus/orca/cli"
	pb "github.com/predixus/orca/protobufs/go"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/reflection"
)

type (
	orcaCoreServer struct {
		pb.UnimplementedOrcaCoreServer
	}
)

var (
	MAX_PROCESSORS = 20
	processors     = make(
		[]grpc.ServerStreamingServer[pb.ProcessingTask],
		MAX_PROCESSORS,
		MAX_PROCESSORS,
	)
)

// CLI
type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	cli.Run()
}

// Register a processor with orca-core. Called on processor startup.
func (orcaCoreServer) RegisterProcessor(
	reg *pb.ProcessorRegistration,
	stream grpc.ServerStreamingServer[pb.ProcessingTask],
) error {
	slog.Info("registering processor",
		"runtime", reg.Runtime)

	// do stuff

	return nil
}

func (orcaCoreServer) EmitWindow(
	ctx context.Context,
	window *pb.Window,
) (*pb.WindowEmitStatus, error) {
	slog.Info("received window",
		"name", window.Name,
		"from", window.From,
		"to", window.To)
	return &pb.WindowEmitStatus{
		Status: pb.WindowEmitStatus_NO_TRIGGERED_ALGORITHMS,
	}, nil
}

func (orcaCoreServer) RegisterWindowType(
	ctx context.Context,
	windowType *pb.WindowType,
) (*pb.Status, error) {
	slog.Info("registering window type",
		"name", windowType.Name)
	return &pb.Status{
		Received: true,
	}, nil
}

func (orcaCoreServer) RegisterAlgorithm(
	ctx context.Context,
	algorithm *pb.Algorithm,
) (*pb.Status, error) {
	slog.Info("registering algorithm",
		"name", algorithm.Name,
		"version", algorithm.Version)
	return &pb.Status{
		Received: true,
	}, nil
}

func (orcaCoreServer) SubmitResult(
	ctx context.Context,
	result *pb.Result,
) (*pb.Status, error) {
	slog.Info("received result",
		"algorithm", result.AlgorithmName,
		"version", result.Version,
		"status", result.Status)
	return &pb.Status{
		Received: true,
	}, nil
}

func (orcaCoreServer) GetDagState(
	ctx context.Context,
	request *pb.DagStateRequest,
) (*pb.DagState, error) {
	slog.Info("getting DAG state",
		"window_id", request.WindowId)
	return &pb.DagState{}, nil
}

func newServer() *orcaCoreServer {
	s := &orcaCoreServer{}
	return s
}

// func main() {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	}))
// 	slog.SetDefault(logger)
//
// 	port := 4040
// 	slog.Debug("Running the server", "port", port)
// 	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
// 	if err != nil {
// 		slog.Error("failed to listen", "message", err)
// 	}
// 	var opts []grpc.ServerOption
// 	grpcServer := grpc.NewServer(opts...)
// 	pb.RegisterOrcaCoreServer(grpcServer, newServer())
// 	reflection.Register(grpcServer)
// 	err = grpcServer.Serve(lis)
// 	if err != nil {
// 		slog.Error("failed to serve", "error", err)
// 	}
// }
