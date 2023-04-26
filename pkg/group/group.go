package group

import (
	"fmt"
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/source"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/samber/lo"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	MaxSelectHeight    = 500
	PreviousGroupAlias = "-"
	SortAsc            = "asc"
	SortDesc           = "desc"
)

type Client struct {
	Config    *config.Config
	State     *state.State
	APIConfig *api.Config
}

func New() (*Client, error) {
	configClient := &config.Client{
		File: config.DefaultConfigPath,
	}
	config, err := configClient.Read()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(config.Global.Kubeconfig)
	if err != nil {
		return nil, err
	}

	state, err := state.Read(config)
	if err != nil {
		return nil, err
	}

	apiConfig, err := kubeconfig.Read(file)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:    config,
		State:     state,
		APIConfig: apiConfig,
	}, nil
}

func (c *Client) Get(groupName string) (*config.GroupItem, error) {
	match, ok := lo.Find(c.Config.Group.Items, func(item config.GroupItem) bool {
		return item.Name == groupName
	})
	if !ok {
		return nil, fmt.Errorf("could not find group: '%s'", groupName)
	}

	return &match, nil
}

func (c *Client) Set(groupName string) error {
	log := logger.New()
	history := c.State.Group.History

	if len(history) > 1 && groupName == PreviousGroupAlias {
		groupName = string(history[len(history)-2])
	}

	if len(groupName) == 0 {
		printer, err := c.buildInteractiveSelectPrinter()
		if err != nil {
			return err
		}

		groupName, err = printer.Show()
		if err != nil {
			return err
		}
	}

	var files []*os.File

	group, ok := lo.Find(c.Config.Group.Items, func(item config.GroupItem) bool {
		return item.Name == groupName
	})
	if !ok {
		return fmt.Errorf("could not find group: '%s", groupName)
	}

	for _, sourceName := range group.Sources {
		sourceMatch, ok := lo.Find(c.Config.Source.Items, func(item config.SourceItem) bool {
			return sourceName == item.Name
		})
		if !ok {
			log.Warn("could not find source", log.Args("source", sourceName, "group", groupName))
			continue
		}
		match, err := source.ComputeFiles(&sourceMatch)
		if err != nil {
			return err
		}
		files = append(files, match...)
	}

	apiConfig, err := kubeconfig.Merge(files...)
	if err != nil {
		return err
	}

	// if the group has a default context, set it
	defaultContext := group.Context.Default
	if len(defaultContext) > 0 {
		contextClient := context.Client{
			Config:    c.Config,
			State:     c.State,
			APIConfig: apiConfig,
		}
		err := contextClient.Set(defaultContext)
		if err != nil {
			return err
		}
	}

	// set new api config and modify state
	c.APIConfig = apiConfig
	c.State.Group.Active = groupName
	c.State.Group.History = state.ComputeHistory(c.Config, state.History(groupName), c.State.Group.History)

	return nil
}

func (c *Client) Reload() error {
	groupName := c.State.Group.Active

	err := c.Set(groupName)
	if err != nil {
		return err
	}
	return nil
}
