package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"slices"
	"strings"
)

type cliFlags struct {
	platform string
	connStr  string
	port     int
	logLevel string
	migrate  bool
	showHelp bool
}

var logLevels = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
}

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

func ValidateLogLevel(s string) error {
	if s == "" {
		return errors.New("You must select a log level")
	}

	s = strings.ToUpper(s)
	for _, level := range logLevels {
		if s == level {
			return nil
		}
	}
	return fmt.Errorf("Invalid log level: %s. Must be one of: %s", s, strings.Join(logLevels, ", "))
}

func parseFlags() cliFlags {
	flags := cliFlags{}

	// connection string
	flag.StringVar(
		&flags.platform,
		"platform",
		"",
		"Data platform to use as the data layer (e.g., postgresql)",
	)
	flag.StringVar(&flags.connStr, "connStr", "", "Connection string to the datalayer")
	flag.IntVar(&flags.port, "port", 4040, "Port number for the Orca server")
	flag.StringVar(&flags.logLevel, "logLevel", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	flag.BoolVar(&flags.showHelp, "help", false, "Show help")
	flag.BoolVar(
		&flags.migrate,
		"migrate",
		false,
		"Migrate the orca db prior to launching orca. Will need to be run at least once to provision the store before use",
	)
	flag.Parse()

	return flags
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func validateFlags(flags cliFlags) error {
	if flags.showHelp {
		return nil
	}

	if flags.platform == "" {
		return fmt.Errorf("a platform selection is required")
	}
	if err := ValidateDatalayer(flags.platform); err != nil {
		return fmt.Errorf("invalid platform: %w", err)
	}

	if flags.connStr == "" {
		return fmt.Errorf("connStr is required")
	}
	if err := ValidateConnStr(flags.connStr); err != nil {
		return fmt.Errorf("invalid connection string: %w", err)
	}

	if err := ValidatePort(fmt.Sprintf("%d", flags.port)); err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}

	if err := ValidateLogLevel(flags.logLevel); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	return nil
}

func runCLI(flags cliFlags) {
	if flags.showHelp {
		flag.Usage()
		return
	}

	// stdout logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(flags.logLevel),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// perform migrations if requested
	slog.Info("premigration")
	if flags.migrate {
		slog.Info("migrating datalayer")
		err := migrateDatalayer(flags.platform, flags.connStr)
		if err != nil {
			slog.Error("could not migrate the datalayer, exiting", "error", err)
			os.Exit(1)
		}
	}
	startGRPCServer(flags.platform, flags.connStr, flags.port, flags.logLevel)

	// Keep main thread alive
	select {}
}
