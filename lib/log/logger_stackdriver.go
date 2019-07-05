package log

import (
	stdlog "log"
	"os"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

// StackDriverLogger outputs logs in a StackDriver compatible format.
type StackDriverLogger struct {
	logger *logrus.Logger

	fields map[string]interface{}
}

// NewStackDriverLogger returns a new StackDriverLogger instance.
func NewStackDriverLogger(cfg *Config) *StackDriverLogger {
	l := logrus.New()

	l.SetFormatter(joonix.NewFormatter(joonix.StackdriverFormat))
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	l.SetLevel(level)

	stdlog.SetOutput(l.Writer())

	return &StackDriverLogger{
		logger: l,

		fields: make(map[string]interface{}),
	}
}

// Debug sends a debug level message.
func (l *StackDriverLogger) Debug(msg string, fields ...Fields) {
	l.log(logrus.DebugLevel, msg, fields...)
}

// Info sends an info level message.
func (l *StackDriverLogger) Info(msg string, fields ...Fields) {
	l.log(logrus.InfoLevel, msg, fields...)
}

// Warn sends a warn level message.
func (l *StackDriverLogger) Warn(msg string, fields ...Fields) {
	l.log(logrus.WarnLevel, msg, fields...)
}

// Error sends an error level message.
func (l *StackDriverLogger) Error(msg string, fields ...Fields) {
	l.log(logrus.ErrorLevel, msg, fields...)
}

// ErrorErr sends an error level message and an associated error.
func (l *StackDriverLogger) ErrorErr(err error, msg string, fields ...Fields) {
	fields = append(fields, Fields{
		"error": err,
	})

	l.log(logrus.ErrorLevel, msg, fields...)
}

// Fatal sends a fatal level message. Terminates execution.
func (l *StackDriverLogger) Fatal(msg string, fields ...Fields) {
	l.log(logrus.FatalLevel, msg, fields...)
	os.Exit(1)
}

// FatalErr sends a fatal level message and an associated error. Terminates execution.
func (l *StackDriverLogger) FatalErr(err error, msg string, fields ...Fields) {
	fields = append(fields, Fields{
		"error": err,
	})

	l.log(logrus.FatalLevel, msg, fields...)
	os.Exit(1)
}

func (l *StackDriverLogger) log(level logrus.Level, msg string, fields ...Fields) {
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
func (l *StackDriverLogger) With(fields Fields) Logger {
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
