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
		Short: "get [context|group] [name], defaults to context",
	}

	cmd.AddCommand(newGroupCommand())
	cmd.AddCommand(newContextCommand())
	return cmd
}
