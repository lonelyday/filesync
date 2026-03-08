package config

import (
	"fmt"
	"os"
	"time"
)

// Logger provides structured logging for filesync operations.
type Logger struct {
	verbose bool
}

// NewLogger creates a new Logger instance.
func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Info logs an informational message.
func (l *Logger) Info(msg string) {
	fmt.Printf("[%s] INFO  %s\n", timestamp(), msg)
}

// Infof logs a formatted informational message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Success logs a success message.
func (l *Logger) Success(msg string) {
	fmt.Printf("[%s] OK    %s\n", timestamp(), msg)
}

// Successf logs a formatted success message.
func (l *Logger) Successf(format string, args ...interface{}) {
	l.Success(fmt.Sprintf(format, args...))
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	fmt.Printf("[%s] WARN  %s\n", timestamp(), msg)
}

// Warnf logs a formatted warning message.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Error logs an error message to stderr.
func (l *Logger) Error(msg string) {
	fmt.Fprintf(os.Stderr, "[%s] ERROR %s\n", timestamp(), msg)
}

// Errorf logs a formatted error message to stderr.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Debug logs a debug message (only when verbose mode is enabled).
func (l *Logger) Debug(msg string) {
	if l.verbose {
		fmt.Printf("[%s] DEBUG %s\n", timestamp(), msg)
	}
}

// Debugf logs a formatted debug message (only when verbose mode is enabled).
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.verbose {
		l.Debug(fmt.Sprintf(format, args...))
	}
}
