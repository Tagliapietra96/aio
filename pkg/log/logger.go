// log package provides a simple logger for the application
package log

import (
	"aio/pkg/utils/fs"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
)

// getLogDir function returns the log directory
func getLogDir() (string, error) {
	execDir, err := fs.ExecDir()
	if err != nil {
		return "", err
	}

	logDir := filepath.Join(execDir, "logs")
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return "", err
	}

	return logDir, nil
}

// loggerInit function initializes the logger
func loggerInit() (*log.Logger, *os.File, error) {
	// get the log directory
	logDir, err := getLogDir()
	if err != nil {
		return nil, nil, err
	}

	// create a log file with the current date if it does not exist
	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(logFile)
	logger.SetReportTimestamp(true)
	logger.SetTimeFormat("[Monday, 02 Jan 2006 15:04:05]")
	logger.SetLevel(log.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetCallerOffset(1)

	return logger, logFile, nil
}

// Deb function logs a debug message
func Deb(msg string, args ...any) {
	logger, file, err := loggerInit() // initialize the logger
	if err != nil {
		Fat(errors.Join(err, file.Close())) // close the file if an error occurs
	}

	defer file.Close()         // close the file when the function ends
	logger.Debug(msg, args...) // log the message
}

// Info function logs an info message
func Info(msg string, args ...any) {
	logger, file, err := loggerInit() // initialize the logger
	if err != nil {
		Fat(errors.Join(err, file.Close())) // close the file if an error occurs
	}

	defer file.Close()        // close the file when the function ends
	logger.Info(msg, args...) // log the message
}

// Warn function logs a warning message
func Warn(msg string, args ...any) {
	logger, file, err := loggerInit() // initialize the logger
	if err != nil {
		Fat(errors.Join(err, file.Close())) // close the file if an error occurs
	}

	defer file.Close()        // close the file when the function ends
	logger.Warn(msg, args...) // log the message
}

// Err function logs an error message
func Err(msg string, args ...any) {
	logger, file, err := loggerInit() // initialize the logger
	if err != nil {
		Fat(errors.Join(err, file.Close())) // close the file if an error occurs
	}

	defer file.Close()         // close the file when the function ends
	logger.Error(msg, args...) // log the message
}

// Fat function logs a fatal error message and exits the program
// it also displays an alert with the error message
func Fat(err error) {
	if err == nil {
		return
	}

	logger, file, err := loggerInit() // initialize the logger
	if err != nil {
		err = errors.Join(err, file.Close()) // close the file if an error occurs
		err = errors.Join(errors.New("failed to initialize logger"), err)
		beeep.Alert("aio: an error occurred", err.Error(), "")
		os.Exit(1)
	}

	logDir, err := getLogDir() // get the log directory
	if err != nil {
		err = errors.Join(err, file.Close()) // close the file if an error occurs
		err = errors.Join(errors.New("failed to get log directory"), err)
		beeep.Alert("aio: an error occurred", err.Error(), "")
		os.Exit(1)
	}

	logger.Error("FATAL", "error", err) // log the message
	prefix := ErrorStyle.Render("ERROR: ")
	println(prefix + "check logs for more details")                                                                // print an advice to check the logs
	beeep.Alert("aio: an error occurred", "To see the full error, check the log file in this folder: "+logDir, "") // display an alert to check the logs
	file.Close()                                                                                                   // close the file
	os.Exit(1)                                                                                                     // exit the program
}
