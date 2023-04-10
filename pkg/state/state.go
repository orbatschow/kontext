package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
)

type History string

type Group struct {
	Active  string    `json:"active,omitempty"`
	History []History `json:"history,omitempty"`
}

type Context struct {
	Active  string    `json:"active,omitempty"`
	History []History `json:"history,omitempty"`
}

type State struct {
	Group   Group   `json:"group"`
	Context Context `json:"context"`
}

var (
	defaultStateDirectory = path.Join(xdg.StateHome, "kontext")
	defaultStateFile      = path.Join(defaultStateDirectory, "state.json")
)

const DefaultMaximumHistorySize = 10

func computeStateFileLocation(config *config.Config) string {
	var stateFile string
	if len(config.State.Location) == 0 {
		stateFile = defaultStateFile
	} else {
		stateFile = config.State.Location
	}

	return stateFile
}

// Init checks if the state directory exists and creates all directories and files if necessary
func Init(config *config.Config) error {
	log := logger.New()
	stateFile := computeStateFileLocation(config)

	// return if the state file already exists
	if _, err := os.Stat(stateFile); err == nil {
		return nil
	}

	log.Debug("missing state file, creating now", log.Args("path", stateFile))

	// create state directory
	baseStateDirectory, _ := filepath.Split(stateFile)
	err := os.MkdirAll(baseStateDirectory, 0755)
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

// Read reads the current state file and serialize it with koanf
func Read(config *config.Config) (*State, error) {
	instance := koanf.New(".")

	log := logger.New()
	stateFile := computeStateFileLocation(config)
	var state *State

	// load the state file into koanf
	if err := instance.Load(file.Provider(stateFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config file, expected file at '%s'", stateFile)
	}

	// unmarshal the state file into struct
	if err := instance.UnmarshalWithConf("", &state, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return nil, fmt.Errorf("could not unmarshal state, err: '%w'", err)
	}
	log.Debug("read state file", log.Args("path", stateFile))

	return state, nil
}

// Write serializes the current state with koanf
func Write(config *config.Config, state *State) error {
	log := logger.New()
	stateFile := computeStateFileLocation(config)

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

// ComputeHistory takes the current history and appends a new entry
// If the history size is larger than the configured or default size, it will remove
// the oldest entry from the history
func ComputeHistory(config *config.Config, entry History, history []History) []History {
	var maxHistorySize int

	if config.State.History.Size == nil {
		maxHistorySize = DefaultMaximumHistorySize
	} else {
		maxHistorySize = *config.State.History.Size
	}

	// if latest entry in history is already equal to the new entry, just return the history
	if len(history) > 0 && history[len(history)-1] == entry {
		return history
	}
	history = append(history, entry)
	if len(history) > maxHistorySize {
		_, history = history[0], history[1:]
	}

	return history
}
