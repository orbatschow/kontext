package cmd

import (
	"os"

	"github.com/orbatschow/kontext/pkg/cmd/get"
	"github.com/orbatschow/kontext/pkg/cmd/reload"
	"github.com/orbatschow/kontext/pkg/cmd/set"
	kubectx "github.com/orbatschow/kontext/pkg/context"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kontext",
	Short: "manage kubernetes config files, contexts, groups and sources",
	Run: func(cmd *cobra.Command, args []string) {
		var name string
		if len(args) == 0 {
			name = ""
		}
		err := kubectx.Set(name)
		if err != nil {
			pterm.Printfln("%v", err)
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
