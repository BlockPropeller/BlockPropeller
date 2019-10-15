package config

import (
	"os"

	"blockpropeller.dev/lib/log"
)

// Config is an abstraction over a configuration structure
// that can be used as the declaration for configuration
// parameters.
//
// The Config structure should mirror how you want the configuration
// laid out in the YAML file.
//
// For example, if you want two configuration options named `foo` and `bar`
// at the top level, you would define the configuration struct like this:
//
// type Config struct {
//   Foo string `yaml:"foo"`
//   Bar string `yaml:"bar"`
// }
//
// Alternatively, you could define an arbitrary nested struct, and those
// properties would appear in the same way in the configuration file.
//
//
// Additionally, with the provided Validate method, the base
// configuration struct, as well as any child struct implementing
// the interface can validate the configured values before returning
// the configuration to the caller.
//
// Keep in mind that in order for Validate to be able to modify config
// structures with default values, said structures must be defined with
// a pointer. The config package internally uses reflection to initialize
// all the pointer structs, so you don't have to do it manually.
//
// @TODO: Document configuration override and extract it into a configurable option.
type Config interface {
	Validate() error
}

// Provider interface is used when initializing project level configuration.
type Provider interface {
	Load(cfg Config) (string, error)
}

// MustLoad matches the signature of `Load` but failing in case of
// unsuccessful configuration load.
func MustLoad(name string, cfg Config, opts ...FileProviderOpt) {
	err := Load(name, cfg, opts...)
	if err != nil {
		log.ErrorErr(err, "could not load service config")

		os.Exit(1)
	}
}

// Load initializes a configuration structure, based on the options provided.
func Load(name string, cfg Config, opts ...FileProviderOpt) error {
	opts = append([]FileProviderOpt{
		WithOverride(name),
		WithName("config"),
		WithPath("config"),
		WithPath("."),
	}, opts...)

	_, err := NewFileProvider(opts...).Load(cfg)

	return err
}
