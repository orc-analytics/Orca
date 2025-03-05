package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type cliFlags struct {
	platform string
	connStr  string
	port     int
	logLevel string
	showHelp bool
}

func parseFlags() (cliFlags, bool) {
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
	flag.IntVar(&flags.port, "P", 4040, "Port number for the Orca server")

	flag.StringVar(&flags.logLevel, "logLevel", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")
	flag.StringVar(&flags.logLevel, "l", "INFO", "Log level (DEBUG, INFO, WARN, ERROR)")

	flag.BoolVar(&flags.showHelp, "help", false, "Show help")
	flag.BoolVar(&flags.showHelp, "h", false, "Show help")

	flag.Parse()

	// Check if any flags were provided
	hasFlags := false
	flag.Visit(func(f *flag.Flag) {
		hasFlags = true
	})

	return flags, hasFlags || flags.showHelp
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

	// Setup stdout logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(flags.logLevel),
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Start server without TUI
	startGRPCServer(flags.platform, flags.connStr, flags.port, flags.logLevel, nil)

	// Keep main thread alive
	select {}
}
