package config_test

import (
	"testing"

	"chainup.dev/lib/config"
)

type TestConfig struct {
}

func (cfg TestConfig) Validate() error {
	return nil
}

func TestFileConfigProvider(t *testing.T) {
	type FooConfig struct {
		TestConfig

		Foo string `yaml:"foo"`
	}

	provider := config.NewFileProvider(config.WithPath("testdata"))

	var cfg FooConfig
	if _, err := provider.Load(&cfg); err != nil {
		t.Errorf("failed loading config: %s", err)
		return
	}

	got := cfg.Foo
	want := "bar"
	if got != want {
		t.Errorf("load config from file: got '%s', want '%s'", got, want)
	}
}

func TestSearchingForPath(t *testing.T) {
	type EnvConfig struct {
		TestConfig
		Env string `yaml:"env"`
	}

	provider := config.NewFileProvider(
		config.WithName("config.yaml.example"),
		config.SearchForPath(),
	)

	var cfg EnvConfig
	if _, err := provider.Load(&cfg); err != nil {
		t.Errorf("failed loading config: %s", err)
		return
	}

	got := cfg.Env
	want := "development"
	if got != want {
		t.Errorf("load config from file: got '%s', want '%s'", got, want)
	}
}

func TestConfigOverriding(t *testing.T) {
	type MySQLConfig struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Secure   bool   `yaml:"secure"`
	}
	type AppConfig struct {
		TestConfig

		MySQL MySQLConfig `yaml:"mysql"`
	}

	provider := config.NewFileProvider(
		config.WithOverride("override"),
		config.WithName("override_config"),
		config.WithPath("testdata"),
	)

	var cfg AppConfig
	if _, err := provider.Load(&cfg); err != nil {
		t.Errorf("failed loading config: %s", err)
		return
	}

	if wantUser, gotUser := "admin", cfg.MySQL.Username; wantUser != gotUser {
		t.Errorf("unexpected username: got %s, want %s", gotUser, wantUser)
	}
	if wantPass, gotPass := "admin", cfg.MySQL.Password; wantPass != gotPass {
		t.Errorf("unexpected password: got %s, want %s", gotPass, wantPass)
	}
	if gotSecure := cfg.MySQL.Secure; gotSecure != true {
		t.Errorf("unexpected secure: got %t, want %t", gotSecure, true)
	}
}

type BazConfig struct {
	Message string
}

func (cfg *BazConfig) Validate() error {
	if cfg.Message == "" {
		cfg.Message = "Hello World"
	}

	return nil
}

type BarConfig struct {
	Baz *BazConfig

	Something *string
}

type FooConfig struct {
	Bar    *BarConfig
	FooBaz BazConfig
}

func (cfg *FooConfig) Validate() error {
	return nil
}

func TestValidation(t *testing.T) {
	var cfg FooConfig

	config.MustLoad("test-service", &cfg, config.SearchForPath())

	if cfg.Bar.Baz.Message != "Hello World" {
		t.Errorf("failed loading config with defaults")
	}
}
