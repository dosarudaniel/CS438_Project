package logger

import (
	"fmt"
	"os"
	"time"
)

type logLevel uint8

/*
	TraceColor   = "\033[0;36m%s\033[0m"
	DebugColor   = "\033[0;96m%s\033[0m"
	InfoColor    = "\033[1;34m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
*/
const (
	// TraceLevel is the most verbose level (TraceLevel => Trace+Debug+Info+Warning+Fatal)
	TraceLevel logLevel = iota
	// DebugLevel is the base debug level (Debug => Debug+Info+Warning+Fatal)
	DebugLevel logLevel = iota
	// InfoLevel is the base logging level (Info => Info+Warning+Fatal)
	InfoLevel logLevel = iota
	// WarningLevel is the first problem level (Warning => Warning+Fatal)
	WarningLevel logLevel = iota
	// FatalLevel is the less verbose level (Fatal => only Fatal)
	FatalLevel logLevel = iota
)

// Logger represents the type of logging we want
type Logger struct {
	Level logLevel
	Timed bool
}

// DefaultLogger creates a new Logger object with default values
func DefaultLogger() *Logger {
	return &Logger{
		Level: FatalLevel,
		Timed: false,
	}
}

func (logger *Logger) timestamp() string {
	if logger.Timed {
		dt := time.Now()
		return fmt.Sprint(" ", dt.Format("2006-01-02 15:04:05"), " ")
	}
	return ""
}

// Trace logs the given string if the logger is enabled for trace outputs
func (logger *Logger) Trace(msg string) {
	if logger.Level <= TraceLevel {
		fmt.Println(logger.timestamp(), "[\033[0;36mTRACE\033[0m]", msg)
	}
}

// Debug logs the given string if the logger is enabled for debug outputs
func (logger *Logger) Debug(msg string) {
	if logger.Level <= DebugLevel {
		fmt.Println(logger.timestamp(), "[\033[0;96mDEBUG\033[0m]", msg)
	}
}

// Info logs the given string if the logger is enabled for normal log outputs
func (logger *Logger) Info(msg string) {
	if logger.Level <= InfoLevel {
		fmt.Println(logger.timestamp(), "[\033[1;34mINFO\033[0m]", msg)
	}
}

// Warn logs the given string if the logger is enabled for warning outputs
func (logger *Logger) Warn(msg string) {
	if logger.Level <= WarningLevel {
		_, err := fmt.Fprintln(os.Stderr, logger.timestamp(), "[\033[1;33mWARN\033[0m]", msg)
		if err != nil {
			fmt.Printf("[\033[1;31mError while warn-logging to stderr:\033[0m\n\tMessage: %v\n\tError: %v\n", msg, err)
		}
	}
}

// Fatal logs the given string if the logger is enabled for fatal outputs
func (logger *Logger) Fatal(msg string) {
	if logger.Level <= FatalLevel {
		_, err := fmt.Fprintln(os.Stderr, logger.timestamp(), "[\033[1;31mFATAL\033[0m]", msg)
		if err != nil {
			fmt.Printf("[\033[1;31mError while fatal-logging to stderr:\033[0m\n\tMessage: %v\n\tError: %v\n", msg, err)
		}
	}
}
