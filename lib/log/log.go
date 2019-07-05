package log

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Config exposes the configuration parameters of the log package
// in the form compatible with the config package.
type Config struct {
	Level string `json:"level" yaml:"level"`
}

// Validate conforms to the config.Config interface,
// providing default values and validating user configuration.
func (cfg *Config) Validate() error {
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	if _, err := logrus.ParseLevel(cfg.Level); err != nil {
		return errors.Wrap(err, "parse log level")
	}

	return nil
}

// Fields is a map holding arbitrary key value data
// to be logged along the message as tags.
type Fields map[string]interface{}

// Logger provides an abstraction over an arbitrary logging library.
//
// This logger includes the capability to do both
// structured and leveled logging.
type Logger interface {
	// Debug sends a debug level message.
	Debug(msg string, fields ...Fields)
	// Info sends an info level message.
	Info(msg string, fields ...Fields)
	// Warn sends a warn level message.
	Warn(msg string, fields ...Fields)
	// Error sends an error level message.
	Error(msg string, fields ...Fields)
	// ErrorErr sends an error level message and an associated error.
	ErrorErr(err error, msg string, fields ...Fields)
	// Fatal sends a fatal level message. Terminates execution.
	Fatal(msg string, fields ...Fields)
	// FatalErr sends a fatal level message and an associated error. Terminates execution.
	FatalErr(err error, msg string, fields ...Fields)

	// With returns a new instance of Logger with the provided Fields attached.
	With(fields Fields) Logger
}
