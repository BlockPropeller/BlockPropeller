package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	// ErrNotFound is an error returned when the configuration does not exist.
	ErrNotFound = errors.New("config not found")
)

// FileProviderOpt allows for specifying custom FileProvider
// configuration while keeping a simple interface.
type FileProviderOpt func(*FileProvider)

// WithOverride option enables overriding the base configuration
// with configuration under the provided key.
//
// Example:
//   foo: foo
//
//   override:
//     foo: bar
// would return `foo = bar` with this option set to `override`
func WithOverride(name string) FileProviderOpt {
	return func(p *FileProvider) {
		p.override = name
	}
}

// WithName configures the name of the configuration file.
//
// By default, `config` is used.
func WithName(name string) FileProviderOpt {
	return func(p *FileProvider) {
		p.configName = name
	}
}

// WithPath adds a new path to search the config in.
//
// This allows for either relative or absolute config paths.
func WithPath(path string) FileProviderOpt {
	return func(p *FileProvider) {
		p.configPaths = append(p.configPaths, path)
	}
}

// SearchForPath looks for a config file in parent folders, iteratively,
// until a config file is found. If not found does nothing.
//
// Set the configuration name before searching for the path
// if using custom config name.
func SearchForPath() FileProviderOpt {
	return func(p *FileProvider) {
		dir, err := os.Getwd()
		if err != nil {
			return
		}

		for dir != "" {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				return
			}

			for _, f := range files {
				if f.IsDir() {
					continue
				}

				if matchesConfig(f.Name(), p) {
					p.configPaths = append(p.configPaths, dir)
					return
				}
			}

			dir, _ = filepath.Split(strings.TrimRight(dir, string(filepath.Separator)))
		}
	}
}

// FileProvider implements the config.Provider interface
// and uses the filesystem as the configuration source.
//
// This is the default configuration provider.
type FileProvider struct {
	configName  string
	configPaths []string
	extensions  []string

	override string
}

// NewFileProvider returns a new FileProvider instance.
func NewFileProvider(opts ...FileProviderOpt) *FileProvider {
	fp := &FileProvider{
		configName:  "config",
		configPaths: []string{},
		extensions:  []string{"", ".yml", ".yaml"},
	}
	for _, opt := range opts {
		opt(fp)
	}

	// If no custom config path is set,
	// the provider looks at the current working directory.
	if len(fp.configPaths) == 0 {
		fp.configPaths = []string{"."}
	}

	return fp
}

// Load satisfies the config.Provider interface.
func (fp *FileProvider) Load(cfg Config) (string, error) {
	path, err := fp.findConfig()
	if err != nil {
		return "", errors.Wrap(err, "find config")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "read config")
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return "", errors.Wrap(err, "parse config")
	}

	if fp.override != "" {
		err = fp.loadOverride(cfg, data)
		if err != nil {
			return "", errors.Wrap(err, "load override config")
		}
	}

	val := reflect.ValueOf(cfg)

	err = fp.initConfig(val)
	if err != nil {
		return "", errors.Wrap(err, "allocate config")
	}

	err = fp.validateConfig(val)
	if err != nil {
		return "", errors.Wrap(err, "validate config")
	}

	return path, nil
}

func (fp *FileProvider) loadOverride(cfg Config, data []byte) error {
	var cfgMap yaml.MapSlice
	err := yaml.Unmarshal(data, &cfgMap)
	if err != nil {
		return errors.Wrap(err, "parse config into slice")
	}
	for _, cfgItem := range cfgMap {
		if cfgItem.Key.(string) != fp.override {
			continue
		}

		var overrideData []byte
		overrideData, err = yaml.Marshal(cfgItem.Value)
		if err != nil {
			return errors.Wrap(err, "marshal override config")
		}

		err = yaml.Unmarshal(overrideData, cfg)
		if err != nil {
			return errors.Wrap(err, "parse override config")
		}
	}

	return nil
}

func (fp *FileProvider) initConfig(val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)

		if f.Kind() == reflect.Ptr && f.Pointer() == 0 &&
			f.Type().Elem().Kind() == reflect.Struct {
			f.Set(reflect.New(f.Type().Elem()))
		}

		if f.Kind() == reflect.Ptr {
			f = reflect.Indirect(f)
		}

		if f.Kind() != reflect.Struct {
			continue
		}

		err := fp.initConfig(f)
		if err != nil {
			return errors.Wrap(err, "allocate field")
		}
	}

	return nil
}

func (fp *FileProvider) validateConfig(val reflect.Value) error {
	if cfg, ok := val.Interface().(Config); ok {
		err := cfg.Validate()
		if err != nil {
			return err
		}
	}

	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)

		if reflect.Indirect(f).Kind() != reflect.Struct {
			continue
		}

		err := fp.validateConfig(f)
		if err != nil {
			return errors.Wrap(err, "validate config")
		}
	}

	return nil
}

func (fp *FileProvider) findConfig() (string, error) {
	for _, path := range fp.configPaths {
		for _, extension := range fp.extensions {
			filePath := filepath.Join(path, fp.configName+extension)
			_, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			return filePath, nil
		}
	}

	return "", ErrNotFound
}

func matchesConfig(name string, p *FileProvider) bool {
	for _, ext := range p.extensions {
		if name == p.configName+ext {
			return true
		}
	}

	return false
}
