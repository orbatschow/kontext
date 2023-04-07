package group

import (
	"fmt"
	"log"
	"strings"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/source"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func Set(groupName string) error {
	log := logger.New()
	kontextConfig := config.Get()
	log.Info("setting group", log.Args("name", groupName))

	var files []string

	group, ok := kontextConfig.Groups[groupName]
	if !ok {
		return fmt.Errorf("could not find group: '%s", groupName)
	}

	for _, sourceRef := range group.Sources {
		value, ok := kontextConfig.Sources[sourceRef]
		if !ok {
			log.Warn("could not find source", log.Args("source", sourceRef, "group", groupName))
			continue
		}
		match, err := source.Expand(&value)
		if err != nil {
			return err
		}
		files = append(files, match...)
	}

	apiConfig, err := kubeconfig.Merge(files...)
	if err != nil {
		return err
	}

	err = kubeconfig.Write(kontextConfig, apiConfig)
	if err != nil {
		return err
	}

	currentState, err := state.Load()
	if err != nil {
		return err
	}

	currentState.GroupState.Active = groupName
	err = state.Write(currentState)
	if err != nil {
		return err
	}

	return nil
}

func Get(cmd *cobra.Command, kontextConfig *config.Config, name string) error {
	buffer := make(map[string]config.Group)

	if len(name) > 0 {
		value, ok := kontextConfig.Groups[name]
		if !ok {
			return fmt.Errorf("could not find group: '%s'", name)
		}
		buffer[name] = value
	} else {
		buffer = kontextConfig.Groups
	}

	err := Print(cmd, buffer)
	if err != nil {
		return err
	}
	return nil
}

func Reload() error {
	currentState, err := state.Load()
	if err != nil {
		return fmt.Errorf("could not load state, err: '%w'", err)
	}
	groupName := currentState.GroupState.Active

	err = Set(groupName)
	if err != nil {
		return err
	}
	return nil
}

func Print(_ *cobra.Command, groups map[string]config.Group) error {
	currentState, err := state.Load()
	if err != nil {
		return err
	}

	table := pterm.TableData{
		{"Active", "Name", "Source(s)"},
	}
	for key, value := range groups {
		active := ""
		if key == currentState.GroupState.Active {
			active = "*"
		}
		table = append(table, []string{
			active, key, strings.Join(value.Sources, "\n"),
		})
	}
	// print empty line for better formatting
	log.Print("")

	// print result table
	err = pterm.DefaultTable.WithHasHeader().WithData(table).Render()
	if err != nil {
		return fmt.Errorf("failed to print table, err: '%w", err)
	}
	return nil
}
