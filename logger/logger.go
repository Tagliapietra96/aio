// logger package provides a simple logger for the application
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	logger  *log.Logger
	logFile *os.File
	logDir  string
	once    sync.Once
)

// loggerInit function initializes the logger
func loggerInit() error {
	var err error
	once.Do(func() {
		// get the path of the executable
		exePath, err := os.Executable()
		if err != nil {
			return
		}

		// make a logs directory in the same directory as the executable if it does not exist
		logDir = filepath.Join(filepath.Dir(exePath), "logs")
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return
		}

		// create a log file with the current date if it does not exist
		logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
		logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		logger = log.New(logFile)
		logger.SetReportTimestamp(true)
		logger.SetTimeFormat("[Monday, 02 Jan 2006 15:04:05]")
		logger.SetLevel(log.DebugLevel)
	})
	return err
}

// Info function logs an info message
func Info(msg string, args ...any) {
	if logger != nil {
		logger.Info(msg, args...)
	}
}

// Debug function logs a debug message
func Debug(msg string, args ...any) {
	if logger != nil {
		logger.Debug(msg, args...)
	}
}

// Warn function logs a warning message
func Warn(msg string, args ...any) {
	if logger != nil {
		logger.Warn(msg, args...)
	}
}

// Error function logs an error message
func Error(msg string, args ...any) {
	if logger != nil {
		logger.Error(msg, args...)
	}
}

// Fatal function logs a fatal message
func Fatal(msg string, err error, args ...any) {
	if err == nil {
		return
	}

	if logger != nil {
		args = append(args, "error", err)
		logger.Error(msg, args...)
		prefix := ErrorStyle.Render("ERROR: ")
		println(prefix + msg)
		println("To see the full error, check the log file in this folder: " + logDir)
		Close()
		os.Exit(1)
	}
}

// Print function prints a message to the console
// with the specified style
func Print(msg string, style lipgloss.Style, args ...any) {
	s := fmt.Sprintf(msg, args...)
	s = style.Render(s)
	println(s)
}

// Line function prints a message to the console
// with a unformatted style
func Line(msg string, args ...any) {
	style := lipgloss.NewStyle()
	Print(msg, style, args...)
}

// Close function closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// init function initializes the package
func init() {
	if err := loggerInit(); err != nil {
		println("Failed to initialize logger:", err)
	}
}
