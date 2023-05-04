package cmd

import (
	"os"

	"github.com/orbatschow/kontext/pkg/cmd/get"
	"github.com/orbatschow/kontext/pkg/cmd/reload"
	"github.com/orbatschow/kontext/pkg/cmd/set"
	"github.com/orbatschow/kontext/pkg/cmd/version"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:    "kontext",
	Short:  "manage kubernetes config files, contexts, groups and sources",
	PreRun: set.Init,
	Run: func(cmd *cobra.Command, args []string) {
		set.NewSetContextCommand(cmd, args)
	},
}

func Execute() {
	// add commands
	rootCmd.AddCommand(get.NewCommand())
	rootCmd.AddCommand(set.NewCommand())
	rootCmd.AddCommand(reload.NewCommand())
	rootCmd.AddCommand(version.NewCommand())

	rootCmd.PersistentFlags().IntVarP(&logger.Verbosity, "verbosity", "v", logger.DefaultVerbosity, "verbose output")

	if err := rootCmd.Execute(); err != nil {
		pterm.Printfln("%v", err)
		os.Exit(1)
	}
}
