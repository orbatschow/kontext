package cmd

import (
	"os"

	"github.com/orbatschow/kontext/pkg/cmd/get"
	"github.com/orbatschow/kontext/pkg/cmd/reload"
	"github.com/orbatschow/kontext/pkg/cmd/set"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kontext",
	Short: "manage kubernetes config files, contexts, groups and sources",
	Run:   set.NewSetContextCommand,
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
