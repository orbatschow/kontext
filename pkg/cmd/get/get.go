package get

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
