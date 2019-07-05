// Package config provides an easy way to statically type your app configuration.
//
// Notable features include:
// - Statically typed configuration.
// - Ability to specify default values.
// - Validation of config parameters.
// - Configuration overrides for multiple services / environments.
//
// Currently, only the FileProvider backend is supported,
// but the package can easily be extended to support other
// sources of configuration, such as environment variables,
// flags and even third party key value stores.
package config
