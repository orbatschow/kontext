package get

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

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
		Use:   "get",
		Short: "Display contexts and groups",
		Long: `List or display details about contexts and groups.

Without arguments, shows available contexts. Use subcommands to list groups or
filter by name.`,
	}

	cmd.AddCommand(newGroupCommand())
	cmd.AddCommand(newContextCommand())
	return cmd
}
