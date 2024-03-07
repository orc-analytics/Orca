package logger

import (
	"github.com/sirupsen/logrus"
)

// an event struct to contain a standard log message
type Event struct {
	message     string
	environment string
	id          int
}

// setup the StandardLogger
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initialises the standard logger
func NewLogger() *StandardLogger {
	baseLogger := logrus.New()

	standardLogger := &StandardLogger{baseLogger}

	standardLogger.Formatter = &logrus.JSONFormatter{}

	return standardLogger
}

// declare standard log formats
var (
	invalidSomethingMessage = Event{"Invalid something: %s", "development", 1}
	Logger                  = NewLogger()
)

// InvalidArg is a standard error message
func (l *StandardLogger) InvalidSomething(argumentName string) {
	l.Errorf(invalidSomethingMessage.message, argumentName)
}
