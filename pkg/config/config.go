package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/pterm/pterm"
)

var (
	DefaultConfigPath = filepath.Join(xdg.ConfigHome, "kontext", "kontext.yaml")
)

const (
	DefaultStateHistoryLimit   = 10
	DefaultBackupRevisionLimit = 10
)

type Client struct {
	// File of the file, that will be used to read the configuration
	File string
}

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
	Directory string `json:"directory,omitempty"`
	Revisions int    `json:"revisions,omitempty"`
}

type History struct {
	Size int `json:"size"`
}

type State struct {
	File    string  `json:"file,omitempty"`
	History History `json:"history,omitempty"`
}

type Global struct {
	Kubeconfig string         `json:"kubeconfig,omitempty"`
	Verbosity  pterm.LogLevel `json:"verbosity,omitempty"`
}

type Config struct {
	Global  Global   `json:"global,omitempty"`
	Backup  Backup   `json:"backup,omitempty"`
	State   State    `json:"state,omitempty"`
	Groups  []Group  `json:"groups,omitempty"`
	Sources []Source `json:"sources,omitempty"`
}

// Read reads the current config file and serialize it with koanf
func (r *Client) Read() (*Config, error) {
	instance := koanf.New(".")
	var config *Config

	// load configuration with into koanf
	if err := instance.Load(file.Provider(r.File), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config file, expected file at '%s'", r.File)
	}

	// set default values
	err := instance.Load(structs.Provider(Config{
		Global: Global{
			Kubeconfig: os.Getenv("KUBECONFIG"),
			Verbosity:  pterm.LogLevelInfo,
		},
		Backup: Backup{
			Enabled:   true,
			Directory: filepath.Join(xdg.DataHome, "kontext", "backup"),
			Revisions: DefaultBackupRevisionLimit,
		},
		State: State{
			History: History{
				Size: DefaultStateHistoryLimit,
			},
			File: filepath.Join(xdg.StateHome, "kontext", "state.json"),
		},
	}, "koanf"), nil)
	if err != nil {
		return nil, err
	}

	// marshal the given configuration into the struct
	if err := instance.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return nil, err
	}

	expandEnvironment(config)

	return config, nil
}

func expandEnvironment(config *Config) {
	config.Global.Kubeconfig = os.ExpandEnv(config.Global.Kubeconfig)
	config.Backup.Directory = os.ExpandEnv(config.Backup.Directory)
	config.State.File = os.ExpandEnv(config.State.File)

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
