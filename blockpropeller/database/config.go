package database

import "github.com/pkg/errors"

// Config object for connecting to a database.
type Config struct {
	Dialect string `yaml:"dialect"`

	// Postgres Options
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Database string `yaml:"database"`

	// Sqlite3 Options
	File string `yaml:"file"`
}

// Validate conforms to the config.Config interface.
func (cfg *Config) Validate() error {
	if cfg.Dialect != "postgres" && cfg.Dialect != "sqlite3" {
		return errors.Errorf(
			"invalid database dialect: '%s'. Valid dialects: 'postgres', 'sqlite3'", cfg.Dialect)
	}

	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 5432
	}
	if cfg.User == "" {
		cfg.User = "postgres"
	}
	if cfg.Pass == "" {
		cfg.Pass = "postgres"
	}
	if cfg.Database == "" {
		cfg.Database = "blockpropeller"
	}

	if cfg.File == "" {
		cfg.File = ".blockpropeller/blockpropeller.db"
	}

	return nil
}
