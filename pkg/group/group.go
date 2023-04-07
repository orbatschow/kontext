package group

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/source"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"k8s.io/client-go/tools/clientcmd/api"
)

const MaxSelectHeight = 500

type Client struct {
	Config    *config.Config
	State     *state.State
	APIConfig *api.Config
}

func New() (*Client, error) {
	config := config.Get()
	file, err := os.Open(config.Global.Kubeconfig)
	if err != nil {
		return nil, err
	}

	apiConfig, err := kubeconfig.Read(file)
	if err != nil {
		return nil, err
	}

	return &Client{
		Config:    config,
		State:     state.Get(),
		APIConfig: apiConfig,
	}, nil
}

func (c *Client) Set(groupName string) error {
	log := logger.New()

	if len(groupName) == 0 {
		var keys []string
		for _, value := range c.Config.Groups {
			keys = append(keys, value.Name)
		}
		groupName, _ = pterm.DefaultInteractiveSelect.WithMaxHeight(MaxSelectHeight).WithOptions(keys).Show()
	}

	var files []string

	group, ok := lo.Find(c.Config.Groups, func(item config.Group) bool {
		return item.Name == groupName
	})
	if !ok {
		return fmt.Errorf("could not find group: '%s", groupName)
	}

	for _, sourceName := range group.Sources {
		sourceMatch, ok := lo.Find(c.Config.Sources, func(item config.Source) bool {
			return sourceName == item.Name
		})
		if !ok {
			log.Warn("could not find source", log.Args("source", sourceName, "group", groupName))
			continue
		}
		match, err := source.Expand(&sourceMatch)
		if err != nil {
			return err
		}
		files = append(files, match...)
	}

	apiConfig, err := kubeconfig.Merge(files...)
	if err != nil {
		return err
	}

	// set new api config and modify state
	c.APIConfig = apiConfig
	c.State.GroupState.Active = groupName

	return nil
}

func (c *Client) Get(groupName string) (*config.Group, error) {
	match, ok := lo.Find(c.Config.Groups, func(item config.Group) bool {
		return item.Name == groupName
	})
	if !ok {
		return nil, fmt.Errorf("could not find group: '%s'", groupName)
	}

	return &match, nil
}

func (c *Client) Reload() error {
	groupName := c.State.GroupState.Active

	err := c.Set(groupName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Print(groups ...config.Group) error {
	table := pterm.TableData{
		{"Active", "Name", "Source(s)"},
	}
	for _, value := range groups {
		active := ""
		if value.Name == c.State.GroupState.Active {
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
