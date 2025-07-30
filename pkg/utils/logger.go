package utils

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger represents a structured logger
type Logger struct {
	verbose bool
	logger  *slog.Logger
}

// NewLogger creates a new logger
func NewLogger(verbose bool) *Logger {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	return &Logger{
		verbose: verbose,
		logger:  slog.New(handler),
	}
}

// Info logs an informational message
func (l *Logger) Info(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...any) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...any) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

// Error logs an error message
func (l *Logger) Error(format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

// Fatal logs a fatal error message and exits
func (l *Logger) Fatal(format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...), slog.String("level", "FATAL"))
	os.Exit(1)
}
