package server

import (
	"time"

	"github.com/pkg/errors"
)

// Config for setting up an HTTP server.
type Config struct {
	// Port to start listening for HTTP requests.
	Port int `yaml:"port"`

	// ReadTimeout in seconds.
	ReadTimeout time.Duration `yaml:"read_timeout"`
	// WriteTimeout in seconds.
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// Validate satisfies the config.Config interface.
func (cfg *Config) Validate() error {
	if cfg.Port == 0 {
		return errors.New("missing server port")
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 30
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 30
	}

	return nil
}
