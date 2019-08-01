package log

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestingLogger provides a user friendly logger while testing.
type TestingLogger struct {
	t *testing.T

	fields map[string]interface{}
}

// NewTestingLogger returns a new TestingLogger instance.
func NewTestingLogger(t *testing.T) *TestingLogger {
	return &TestingLogger{
		t: t,

		fields: make(map[string]interface{}),
	}
}

// Debug sends a debug level message.
func (l *TestingLogger) Debug(msg string, fields ...Fields) {
	l.log(logrus.DebugLevel, msg, fields...)
}

// Info sends an info level message.
func (l *TestingLogger) Info(msg string, fields ...Fields) {
	l.log(logrus.InfoLevel, msg, fields...)
}

// Warn sends a warn level message.
func (l *TestingLogger) Warn(msg string, fields ...Fields) {
	l.log(logrus.WarnLevel, msg, fields...)
}

// Error sends an error level message.
func (l *TestingLogger) Error(msg string, fields ...Fields) {
	l.log(logrus.ErrorLevel, msg, fields...)
}

// ErrorErr sends an error level message and an associated error.
func (l *TestingLogger) ErrorErr(err error, msg string, fields ...Fields) {
	fields = append(fields, Fields{
		"error": err,
	})

	l.log(logrus.ErrorLevel, msg, fields...)
}

// Fatal sends a fatal level message. Terminates execution.
func (l *TestingLogger) Fatal(msg string, fields ...Fields) {
	l.log(logrus.FatalLevel, msg, fields...)
	os.Exit(1)
}

// FatalErr sends a fatal level message and an associated error. Terminates execution.
func (l *TestingLogger) FatalErr(err error, msg string, fields ...Fields) {
	fields = append(fields, Fields{
		"error": err,
	})

	l.log(logrus.FatalLevel, msg, fields...)
	os.Exit(1)
}

func (l *TestingLogger) log(level logrus.Level, msg string, fields ...Fields) {
	var args []string
	for _, f := range fields {
		for key, value := range f {
			args = append(args, fmt.Sprintf("%s=%v", key, value))
		}
	}
	for key, value := range l.fields {
		args = append(args, fmt.Sprintf("%s=%v", key, value))
	}

	l.t.Logf("[%s] %s    %s", strings.ToUpper(level.String()), msg, strings.Join(args, " "))
}

// With returns a new instance of Logger with the provided Fields attached.
func (l *TestingLogger) With(fields Fields) Logger {
	l2 := *l
	l2.fields = make(map[string]interface{})

	for key, value := range l.fields {
		l2.fields[key] = value
	}

	for key, value := range fields {
		l2.fields[key] = value
	}

	return &l2
}
