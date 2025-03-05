package main

import (
	"context"
	"log/slog"
	"time"
)

// TUILogHandler implements slog.Handler interface
type TUILogHandler struct {
	level   slog.Level
	logChan chan logMsg
}

func NewTUILogHandler(logChan chan logMsg, level slog.Level) *TUILogHandler {
	return &TUILogHandler{
		level:   level,
		logChan: logChan,
	}
}

func (h *TUILogHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < h.level {
		return nil
	}

	h.logChan <- logMsg{
		level:   r.Level,
		message: r.Message,
		time:    time.Now(),
	}
	return nil
}

func (h *TUILogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *TUILogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *TUILogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
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
