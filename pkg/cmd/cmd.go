package cmd

import (
	"os"

	"github.com/orbatschow/kontext/pkg/cmd/get"
	"github.com/orbatschow/kontext/pkg/cmd/reload"
	"github.com/orbatschow/kontext/pkg/cmd/set"
	"github.com/orbatschow/kontext/pkg/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:    "kontext",
	Short:  "manage kubernetes config files, contexts, groups and sources",
	PreRun: set.Init,
	Run: func(cmd *cobra.Command, args []string) {
		isVersionFlagSet := cmd.Flags().Lookup("version").Changed
		if isVersionFlagSet {
			pterm.Printfln(version.Compute())
			os.Exit(0)
		}

		set.NewSetContextCommand(cmd, args)
	},
}

func Execute() {
	// add commands
	rootCmd.AddCommand(get.NewCommand())
	rootCmd.AddCommand(set.NewCommand())
	rootCmd.AddCommand(reload.NewCommand())
	rootCmd.Flags().BoolP("version", "v", false, "version for kontext")

	if err := rootCmd.Execute(); err != nil {
		pterm.Printfln("%v", err)
		os.Exit(1)
	}
}
