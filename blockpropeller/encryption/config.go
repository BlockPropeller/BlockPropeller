package encryption

import "github.com/pkg/errors"

type Config struct {
	Secret string `yaml:"secret"`
}

func (cfg *Config) Validate() error {
	if cfg.Secret == "" {
		return errors.New("missing encryption secret key")
	}

	return nil
}
