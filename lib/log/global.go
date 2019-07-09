package log

import (
	"io"
)

var logger Logger

func init() {
	logger = NewConsoleLogger(&Config{})
}

// SetGlobal overrides the default globally accessible logger.
//
// If the log package is primarily used from the global scope,
// SetGlobal can be used to setup some application wide tags
// such as the running service name, application version, etc...
func SetGlobal(l Logger) {
	logger = l
}

// Debug sends a debug level entry to the global Logger instance.
func Debug(msg string, fields ...Fields) {
	logger.Debug(msg, fields...)
}

// Info sends an info level entry to the global Logger instance.
func Info(msg string, fields ...Fields) {
	logger.Info(msg, fields...)
}

// Warn sends a warn level entry to the global Logger instance.
func Warn(msg string, fields ...Fields) {
	logger.Warn(msg, fields...)
}

// Error sends an error level entry to the global Logger instance.
func Error(msg string, fields ...Fields) {
	logger.Error(msg, fields...)
}

// ErrorErr sends an error level entry with an associated error to the global Logger instance.
func ErrorErr(err error, msg string, fields ...Fields) {
	logger.ErrorErr(err, msg, fields...)
}

// Fatal sends a fatal level entry to the global Logger instance.
func Fatal(msg string, fields ...Fields) {
	logger.Fatal(msg, fields...)
}

// FatalErr sends a fatal level entry with an associated error to the global Logger instance.
func FatalErr(err error, msg string, fields ...Fields) {
	logger.FatalErr(err, msg, fields...)
}

// With returns a new Logger instance based on the global Logger, with the provided fields added.
func With(fields Fields) Logger {
	return logger.With(fields)
}

// Closer is a utility function that accepts an io.Closer interface,
// and logs an error in case of the io.Closer failing.
//
// This can be used to wrap an io.Closer called via defer,
// so we can both defer and log in a single line, simplifying code.
//
// Example:
//   db := database.Open()
//   defer log.Closer(db)
func Closer(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		ErrorErr(err, "failed running closer")
	}
}
