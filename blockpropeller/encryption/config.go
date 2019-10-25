package encryption

import "github.com/pkg/errors"

// Config for the encryption module.
type Config struct {
	Secret string `yaml:"secret"`
}

// Validate satisfies the Config interface.
func (cfg *Config) Validate() error {
	if cfg.Secret == "" {
		return errors.New("missing encryption secret key")
	}

	return nil
}
