package context

import (
	"fmt"
	"os"
	"sort"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Client struct {
	Config    *config.Config
	State     *state.State
	APIConfig *api.Config
}

const (
	MaxSelectHeight      = 500
	PreviousContextAlias = "-"
	SortAsc              = "asc"
	SortDesc             = "desc"
)

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

func (c *Client) Get(contextName string) (map[string]*api.Context, error) {
	log := logger.New()
	log.Info("getting context", log.Args("name", contextName))

	if len(contextName) == 0 {
		return nil, fmt.Errorf("given context name is empty")
	}

	buffer, ok := c.APIConfig.Contexts[contextName]
	if !ok {
		return nil, fmt.Errorf("could not find context '%s'", contextName)
	}

	return map[string]*api.Context{
		contextName: buffer,
	}, nil
}

func (c *Client) List() map[string]*api.Context {
	log := logger.New()
	log.Info("listing contexts")

	return c.APIConfig.Contexts
}

func (c *Client) Set(contextName string) error {
	log := logger.New()
	history := c.State.Context.History

	if len(history) > 1 && contextName == PreviousContextAlias {
		contextName = string(history[len(history)-2])
	}

	if len(contextName) == 0 {
		printer, err := c.selectContext()
		if err != nil {
			return err
		}
		contextName, err = printer.Show()
		if err != nil {
			return err
		}
	}

	_, ok := c.APIConfig.Contexts[contextName]
	if !ok {
		return fmt.Errorf("could not find context: '%s'", contextName)
	}

	c.APIConfig.CurrentContext = contextName
	c.State.Context.Active = contextName
	c.State.Context.History = state.ComputeHistory(c.Config, state.History(contextName), history)

	log.Info("switched context", log.Args("context", contextName))
	return nil
}

// start an interactive context selection
func (c *Client) selectContext() (*pterm.InteractiveSelectPrinter, error) {
	// compute all selection options
	var keys []string
	for k := range c.APIConfig.Contexts {
		keys = append(keys, k)
	}

	// get the active group
	group, ok := lo.Find(c.Config.Group.Items, func(item config.GroupItem) bool {
		return item.Name == c.State.Group.Active
	})
	if !ok {
		return nil, fmt.Errorf("could not find active group: '%s'", c.State.Group.Active)
	}

	// sort the selection
	switch group.Context.Selection.Sort {
	case SortAsc:
		sort.Strings(keys)
	case SortDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	default:
		sort.Strings(keys)
	}

	// set the default selection option
	if len(group.Context.Selection.Default) > 0 {
		return pterm.DefaultInteractiveSelect.
			WithMaxHeight(MaxSelectHeight).
			WithOptions(keys).
			WithDefaultOption(group.Context.Selection.Default), nil
	}
	return pterm.DefaultInteractiveSelect.
		WithMaxHeight(MaxSelectHeight).
		WithOptions(keys), nil
}

func (c *Client) Print(contexts map[string]*api.Context) error {
	table := pterm.TableData{
		{"Active", "Name", "Cluster", "AuthInfo"},
	}
	for key, value := range contexts {
		active := ""
		if key == c.State.Context.Active {
			active = "*"
		}
		table = append(table, []string{
			active, key, value.Cluster, value.AuthInfo,
		})
	}
	err := pterm.DefaultTable.WithHasHeader().WithData(table).Render()
	if err != nil {
		return fmt.Errorf("failed to print table, err: '%w", err)
	}
	return nil
}
