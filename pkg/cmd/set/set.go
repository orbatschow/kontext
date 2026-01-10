package set

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

//revive:disable:var-naming
func init_(_ *cobra.Command, _ []string) {
	log := logger.New()
	configClient := &config.Client{
		File: config.DefaultConfigPath,
	}
	currentConfig, err := configClient.Read()
	if err != nil {
		log.Error("failed to read config", "err", err)
		os.Exit(1)
	}

	err = state.Init(currentConfig)
	if err != nil {
		log.Error("failed to initialize state", "err", err)
		os.Exit(1)
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
