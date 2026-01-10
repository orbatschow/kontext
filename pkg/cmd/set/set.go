package set

import (
	"log"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func init_(_ *cobra.Command, _ []string) {
	configClient := &config.Client{
		File: config.DefaultConfigPath,
	}
	currentConfig, err := configClient.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = state.Init(currentConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Switch the active context or group",
		Long: `Change the active Kubernetes context or group.

When called without arguments, opens an interactive selection dialog for the active group's contexts.
Use 'set context' or 'set group' subcommands for more control.`,
	}

	cmd.AddCommand(newGroupCommand())
	cmd.AddCommand(newContextCommand())

	return cmd
}
