package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
)

type Group struct {
	Active  string   `json:"active,omitempty"`
	History []string `json:"history,omitempty"`
}

type Context struct {
	Active  string   `json:"active,omitempty"`
	History []string `json:"history,omitempty"`
}

type State struct {
	Group   Group   `json:"group"`
	Context Context `json:"context"`
}

var (
	instance = koanf.New(".")

	stateDirectory = path.Join(xdg.StateHome, "kontext")
	stateFile      = path.Join(stateDirectory, "state.json")

	state *State
)

const DefaultMaximumHistorySize = 10

func initialize() error {
	log := logger.New()

	// check if the state file already exists
	_, err := os.Stat(stateFile)
	errors.Is(err, os.ErrNotExist)
	if err == nil {
		return nil
	}

	log.Debug("missing state file, creating now", log.Args("path", stateFile))

	// create state directory
	err = os.MkdirAll(stateDirectory, 0755)
	if err != nil {
		return fmt.Errorf("could not create state directory, err: '%w'", err)
	}

	// create state file
	_, err = os.Create(stateFile)
	if err != nil {
		return fmt.Errorf("could not create state file, err: '%w'", err)
	}

	return nil
}

func Read() error {
	log := logger.New()
	err := initialize()
	if err != nil {
		return err
	}

	// load the state file into koanf
	if err := instance.Load(file.Provider(stateFile), yaml.Parser()); err != nil {
		return fmt.Errorf("failed to load config file, expected file at '%s'", stateFile)
	}

	// unmarshal the state file into struct
	if err := instance.UnmarshalWithConf("", &state, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return fmt.Errorf("could not unmarshal state, err: '%w'", err)
	}
	log.Debug("read state file", log.Args("path", stateFile))

	return nil
}

func Write(state *State) error {
	log := logger.New()

	// marshal the state into json
	buffer, err := json.Marshal(state)
	if err != nil {
		return err
	}

	log.Debug("updating state", log.Args("data", string(buffer)))

	// write the state into the state file
	err = os.WriteFile(stateFile, buffer, 0600)
	if err != nil {
		return fmt.Errorf("could not write state to file, err: '%w'", err)
	}
	log.Debug("finished updating state")

	return nil
}

func Get() *State {
	return state
}

func ComputeHistory(config *config.Config, entry string, history []string) []string {
	var maxHistorySize int

	if config.History.Size == nil {
		maxHistorySize = DefaultMaximumHistorySize
	} else {
		maxHistorySize = *config.History.Size
	}

	if len(history) > 0 && history[len(history)-1] == entry {
		return history
	}
	history = append(history, entry)
	if len(history) > maxHistorySize {
		_, history = history[0], history[1:]
	}

	return history
}
