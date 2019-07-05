package log

import (
	stdlog "log"
	"os"

	"github.com/sirupsen/logrus"
)

// ConsoleLogger provides a user friendly logger,
// ideal for the development environment.
type ConsoleLogger struct {
	logger *logrus.Logger

	fields map[string]interface{}
}

// NewConsoleLogger returns a new ConsoleLogger instance.
func NewConsoleLogger(cfg *Config) *ConsoleLogger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{})

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	l.SetLevel(level)

	stdlog.SetOutput(l.Writer())

	return &ConsoleLogger{
		logger: l,

		fields: make(map[string]interface{}),
	}
}

// Debug sends a debug level message.
func (l *ConsoleLogger) Debug(msg string, fields ...Fields) {
	l.log(logrus.DebugLevel, msg, fields...)
}

// Info sends an info level message.
func (l *ConsoleLogger) Info(msg string, fields ...Fields) {
	l.log(logrus.InfoLevel, msg, fields...)
}

// Warn sends a warn level message.
func (l *ConsoleLogger) Warn(msg string, fields ...Fields) {
	l.log(logrus.WarnLevel, msg, fields...)
}

// Error sends an error level message.
func (l *ConsoleLogger) Error(msg string, fields ...Fields) {
	l.log(logrus.ErrorLevel, msg, fields...)
}

// ErrorErr sends an error level message and an associated error.
func (l *ConsoleLogger) ErrorErr(err error, msg string, fields ...Fields) {
	fields = append(fields, Fields{
		"error": err,
	})

	l.log(logrus.ErrorLevel, msg, fields...)
}

// Fatal sends a fatal level message. Terminates execution.
func (l *ConsoleLogger) Fatal(msg string, fields ...Fields) {
	l.log(logrus.FatalLevel, msg, fields...)
	os.Exit(1)
}

// FatalErr sends a fatal level message and an associated error. Terminates execution.
func (l *ConsoleLogger) FatalErr(err error, msg string, fields ...Fields) {
	fields = append(fields, Fields{
		"error": err,
	})

	l.log(logrus.FatalLevel, msg, fields...)
	os.Exit(1)
}

func (l *ConsoleLogger) log(level logrus.Level, msg string, fields ...Fields) {
	args := logrus.Fields{}
	for _, f := range fields {
		for key, value := range f {
			args[key] = value
		}
	}
	for key, value := range l.fields {
		args[key] = value
	}

	l.logger.WithFields(args).Log(level, msg)
}

// With returns a new instance of Logger with the provided Fields attached.
func (l *ConsoleLogger) With(fields Fields) Logger {
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
