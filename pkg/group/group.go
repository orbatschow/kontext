package group

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/source"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
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
		groupName = c.selectGroup()
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

// start an interactive context selection
func (c *Client) selectGroup() string {
	// compute all selection options
	var keys []string
	for _, value := range c.Config.Group.Items {
		keys = append(keys, value.Name)
	}

	// sort the selection
	switch c.Config.Group.Selection.Sort {
	case SortAsc:
		sort.Strings(keys)
	case SortDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	default:
		sort.Strings(keys)
	}

	// set the default selection option
	var groupName string
	if len(c.Config.Group.Selection.Default) > 0 {
		groupName, _ = pterm.DefaultInteractiveSelect.
			WithMaxHeight(MaxSelectHeight).
			WithOptions(keys).
			WithDefaultOption(c.Config.Group.Selection.Default).
			Show()
	} else {
		groupName, _ = pterm.DefaultInteractiveSelect.
			WithMaxHeight(MaxSelectHeight).
			WithOptions(keys).
			Show()
	}

	return groupName
}

func (c *Client) Reload() error {
	groupName := c.State.Group.Active

	err := c.Set(groupName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Print(groups ...config.GroupItem) error {
	table := pterm.TableData{
		{"Active", "Name", "Source(s)"},
	}
	for _, value := range groups {
		active := ""
		if value.Name == c.State.Group.Active {
			active = "*"
		}
		table = append(table, []string{
			active, value.Name, strings.Join(value.Sources, "\n"),
		})
	}
	// print empty line for better formatting
	log.Print("")

	// print result table
	err := pterm.DefaultTable.WithHasHeader().WithData(table).Render()
	if err != nil {
		return fmt.Errorf("failed to print table, err: '%w", err)
	}
	return nil
}
