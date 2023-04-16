package context

import (
	"fmt"
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
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

	if len(history) > 1 && contextName == "-" {
		contextName = string(history[len(history)-2])
	}

	if len(contextName) == 0 {
		var keys []string
		for k := range c.APIConfig.Contexts {
			keys = append(keys, k)
		}
		contextName, _ = pterm.DefaultInteractiveSelect.WithMaxHeight(MaxSelectHeight).WithOptions(keys).Show()
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
