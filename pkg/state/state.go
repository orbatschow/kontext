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
	"github.com/orbatschow/kontext/pkg/logger"
)

type GroupState struct {
	Active string `json:"active,omitempty"`
}

type ContextState struct {
	Active string `json:"active,omitempty"`
}

type State struct {
	GroupState   GroupState   `json:"groupState"`
	ContextState ContextState `json:"contextState"`
}

var (
	instance = koanf.New(".")

	stateDirectory = path.Join(xdg.StateHome, "kontext")
	stateFile      = path.Join(stateDirectory, "state.json")
)

func Init() error {
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

func Load() (*State, error) {
	log := logger.New()
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
