package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
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

type Config struct {
	Global Global `json:"global,omitempty"`
	State  State  `json:"state,omitempty"`
	Backup Backup `json:"backup,omitempty"`
	Group  Group  `json:"group,omitempty"`
	Source Source `json:"source,omitempty"`
}

type Global struct {
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

// State configuration options
type State struct {
	// path of the state file
	File    string  `json:"file,omitempty"`
	History History `json:"history,omitempty"`
}

type History struct {
	// set the maximum history size
	Size int `json:"size"`
}

// Backup configuration options
type Backup struct {
	// enable/disable the backup, defaults to true
	Enabled bool `json:"enabled"`
	// set the backup directory
	Directory string `json:"directory,omitempty"`
	// set the maximum backup revision count
	Revisions int `json:"revisions,omitempty"`
}

// Group configuration Options
type Group struct {
	Items     []GroupItem `json:"items"`
	Selection Selection   `json:"selection"`
}

type GroupItem struct {
	Name    string   `json:"name"`
	Context Context  `json:"context,omitempty"`
	Sources []string `json:"sources"`
}

type Context struct {
	Default   string    `json:"default"`
	Selection Selection `json:"selection"`
}

type Selection struct {
	Default string `json:"default"`
	Sort    string `json:"sort"`
}

// Source configuration options
type Source struct {
	Items []SourceItem `json:"items"`
}

type SourceItem struct {
	Name    string   `json:"name"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude"`
}

// Read reads the current config file and serialize it with koanf
func (r *Client) Read() (*Config, error) {
	instance := koanf.New(".")
	var config *Config

	// load configuration with into koanf
	if err := instance.Load(file.Provider(r.File), yaml.Parser()); err != nil {
		return nil, err
	}

	// set default values
	err := instance.Load(structs.Provider(Config{
		Global: Global{
			Kubeconfig: os.Getenv("KUBECONFIG"),
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

	for i, source := range config.Source.Items {
		for j, include := range source.Include {
			source.Include[j] = os.ExpandEnv(include)
		}
		for j, exclude := range source.Exclude {
			source.Exclude[j] = os.ExpandEnv(exclude)
		}
		config.Source.Items[i] = source
	}
}
