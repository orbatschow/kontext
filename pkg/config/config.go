package config

import (
	"fmt"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/pterm/pterm"
)

// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
var (
	instance = koanf.New(".")
)

var config *Config

type Source struct {
	Name    string   `json:"name"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude"`
}

type Group struct {
	Name    string   `json:"name"`
	Context string   `json:"context,omitempty"`
	Sources []string `json:"sources"`
}

type Backup struct {
	Enabled   bool   `json:"enabled"`
	Location  string `json:"location,omitempty"`
	Revisions *int   `json:"revisions,omitempty"`
}

type History struct {
	Size *int `json:"size"`
}

type State struct {
	Location string  `json:"location,omitempty"`
	History  History `json:"history,omitempty"`
}

type Global struct {
	Kubeconfig string         `json:"kubeconfig,omitempty"`
	Verbosity  pterm.LogLevel `json:"verbosity,omitempty"`
}

type Config struct {
	Global  Global   `json:"global,omitempty"`
	Backup  Backup   `json:"backup"`
	State   State    `json:"state,omitempty"`
	Groups  []Group  `json:"groups"`
	Sources []Source `json:"sources"`
}

// validate will check the given configuration for errors
func validate(config *Config) error {
	if len(config.Global.Kubeconfig) == 0 {
		value, ok := os.LookupEnv("KUBECONFIG")
		if !ok {
			return fmt.Errorf("no kubeconfig path provided and KUBECONFIG environment variable unset")
		}

		config.Global.Kubeconfig = value
	}

	return nil
}

// Read will parse a kontext configuration file
func Read() error {
	configFile := path.Join(xdg.ConfigHome, "kontext", "kontext.yaml")

	if err := instance.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		return fmt.Errorf("failed to load config file, expected file at '%s'", configFile)
	}

	if err := instance.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return err
	}

	err := validate(config)
	if err != nil {
		return err
	}

	expandEnvironment(config)

	return nil
}

func expandEnvironment(config *Config) {
	config.Global.Kubeconfig = os.ExpandEnv(config.Global.Kubeconfig)
	config.Backup.Location = os.ExpandEnv(config.Backup.Location)

	for i, source := range config.Sources {
		for j, include := range source.Include {
			source.Include[j] = os.ExpandEnv(include)
		}
		for j, exclude := range source.Exclude {
			source.Exclude[j] = os.ExpandEnv(exclude)
		}
		config.Sources[i] = source
	}
}

// Get will return a parsed kontext Config struct
func Get() *Config {
	return config
}
