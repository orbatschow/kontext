package cmd

import (
	"os"

	"github.com/orbatschow/kontext/pkg/cmd/get"
	"github.com/orbatschow/kontext/pkg/cmd/reload"
	"github.com/orbatschow/kontext/pkg/cmd/set"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kontext",
	Short: "manage kubernetes config files, contexts, groups and sources",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New()
		var contextName string
		if len(args) == 0 {
			contextName = ""
		}

		client, err := context.New()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		err = client.Set(contextName)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		err = state.Write(client.State)
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	// add commands
	rootCmd.AddCommand(get.NewCommand())
	rootCmd.AddCommand(set.NewCommand())
	rootCmd.AddCommand(reload.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		pterm.Printfln("%v", err)
		os.Exit(1)
	}
}
