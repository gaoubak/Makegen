package utils

import (
	"fmt"
	"os"
	"time"
)

// Logger provides structured logging
type Logger struct {
	verbose bool
}

// NewLogger creates a new logger
func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log("ℹ️", format, args...)
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) {
	l.log("✅", format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log("⚠️", format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log("❌", format, args...)
}

// Debug logs a debug message (only if verbose)
func (l *Logger) Debug(format string, args ...interface{}) {
	if !l.verbose {
		return
	}
	l.log("🔍", format, args...)
}

func (l *Logger) log(emoji string, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] %s %s\n", timestamp, emoji, message)
}
