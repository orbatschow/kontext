package context

import (
	"fmt"
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd/api"
)

type JSONContext struct {
	Name     string `json:"name,omitempty"`
	Cluster  string `json:"cluster,omitempty"`
	AuthInfo string `json:"authInfo,omitempty"`
}

const MaxSelectHeight = 500

func Get(cmd *cobra.Command, contextName string) error {
	log := logger.New()
	kontextConfig := config.Get()
	log.Info("setting context", log.Args("contextName", contextName))

	buffer := make(map[string]*api.Context)

	file, err := os.Open(kontextConfig.Global.Kubeconfig)
	if err != nil {
		return err
	}
	apiConfig, err := kubeconfig.Load(file)
	if err != nil {
		return err
	}

	if len(contextName) > 0 {
		ctx, ok := apiConfig.Contexts[contextName]
		if !ok {
			return fmt.Errorf("could not find context '%s'", contextName)
		}
		buffer[contextName] = ctx
	} else {
		buffer = apiConfig.Contexts
	}

	err = Print(cmd, buffer, apiConfig)
	if err != nil {
		return err
	}
	return nil
}

func Set(contextName string) error {
	log := logger.New()
	kontextConfig := config.Get()

	file, err := os.Open(kontextConfig.Global.Kubeconfig)
	if err != nil {
		return err
	}
	apiConfig, err := kubeconfig.Load(file)
	if err != nil {
		return err
	}

	if len(contextName) == 0 {
		var keys []string
		for k := range apiConfig.Contexts {
			keys = append(keys, k)
		}
		contextName, _ = pterm.DefaultInteractiveSelect.WithMaxHeight(MaxSelectHeight).WithOptions(keys).Show()
	}

	_, ok := apiConfig.Contexts[contextName]
	if !ok {
		return fmt.Errorf("could not find context: '%s'", contextName)
	}

	apiConfig.CurrentContext = contextName
	err = kubeconfig.Write(kontextConfig, apiConfig)
	if err != nil {
		return err
	}

	currentState, err := state.Load()
	if err != nil {
		return err
	}

	currentState.ContextState.Active = contextName
	err = state.Write(currentState)
	if err != nil {
		return err
	}

	log.Info("switched context", log.Args("context", contextName))
	return nil
}

func Reload() error {
	currentState, err := state.Load()
	if err != nil {
		return fmt.Errorf("could not load state, err: '%w'", err)
	}
	contextName := currentState.ContextState.Active
	if len(contextName) == 0 {
		return fmt.Errorf("no active context")
	}

	err = Set(contextName)
	if err != nil {
		return err
	}

	return nil
}

func Print(_ *cobra.Command, contextList map[string]*api.Context, apiConfig *api.Config) error {
	table := pterm.TableData{
		{"Active", "Name", "Cluster", "AuthInfo"},
	}
	for key, value := range contextList {
		active := ""
		if key == apiConfig.CurrentContext {
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
