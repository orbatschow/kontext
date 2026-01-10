package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

var (
	DefaultConfigPath = filepath.Join(xdg.ConfigHome, "kontext", "kontext.yaml")
	validate          = validator.New()
)

const (
	DefaultStateHistoryLimit = 10
)

type Client struct {
	// File of the file, that will be used to read the configuration
	File string
}

type Config struct {
	Global Global `json:"global,omitempty"`
	State  State  `json:"state,omitempty"`
	Groups Group  `json:"groups,omitempty"`
	Sources Source `json:"sources,omitempty"`
}

type Global struct {
	Kubeconfig string `json:"kubeconfig,omitempty"`
}

// State configuration options
type State struct {
	// path of the state file
	File    string  `json:"file,omitempty" validate:"required"`
	History History `json:"history,omitempty"`
}

type History struct {
	// set the maximum history size
	Size int `json:"size" validate:"gte=0"`
}

// Group configuration Options
type Group struct {
	Items     []GroupItem `json:"items,omitempty" validate:"dive"`
	Selection Selection   `json:"selection,omitempty"`
}

type GroupItem struct {
	Name    string   `json:"name" validate:"required,min=1"`
	Context Context  `json:"context,omitempty"`
	Sources []string `json:"sources,omitempty" validate:"omitempty,min=1,dive,min=1"`
}

type Context struct {
	Default   string    `json:"default,omitempty" validate:"required_with=Selection"`
	Selection Selection `json:"selection,omitempty"`
}

type Selection struct {
	Default string `json:"default,omitempty"`
	Sort    string `json:"sort,omitempty" validate:"omitempty,oneof=asc desc"`
}

// Source configuration options
type Source struct {
	Items []SourceItem `json:"items,omitempty" validate:"dive"`
}

type SourceItem struct {
	Name    string   `json:"name" validate:"required,min=1"`
	Include []string `json:"include,omitempty" validate:"omitempty,min=1,dive,min=1"`
	Exclude []string `json:"exclude,omitempty" validate:"omitempty,min=1,dive,min=1"`
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

	// validate the configuration
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	expandEnvironment(config)

	return config, nil
}

func expandEnvironment(config *Config) {
	config.Global.Kubeconfig = os.ExpandEnv(config.Global.Kubeconfig)
	config.State.File = os.ExpandEnv(config.State.File)

	for i, source := range config.Sources.Items {
		for j, include := range source.Include {
			source.Include[j] = os.ExpandEnv(include)
		}
		for j, exclude := range source.Exclude {
			source.Exclude[j] = os.ExpandEnv(exclude)
		}
		config.Sources.Items[i] = source
	}
}
