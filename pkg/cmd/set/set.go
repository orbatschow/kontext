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
		Use:   "set [name]",
		Short: "set [context|group] [name]",
		Long: `Invoking this command without a parameter will spawn an interactive selection dialog.
When providing a context name, the switch will be performed immediately.
'-' is a reserved context name, that will cause a switch to the previously active context.
If neither, context nor group is specified, this command will set the context.
		`,
	}

	cmd.AddCommand(newGroupCommand())
	cmd.AddCommand(newContextCommand())

	return cmd
}
