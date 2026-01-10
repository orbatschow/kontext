package cmd

import (
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/orbatschow/kontext/pkg/cmd/get"
	"github.com/orbatschow/kontext/pkg/cmd/reload"
	"github.com/orbatschow/kontext/pkg/cmd/set"
	"github.com/orbatschow/kontext/pkg/cmd/version"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kontext",
	Short: "Manage Kubernetes kubeconfig files, contexts, and groups",
	Long: `Kontext helps you manage multiple Kubernetes kubeconfig files by organizing them
into groups and providing quick switching between contexts and clusters.`,
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&logger.Verbosity, "verbosity", "v", logger.DefaultVerbosity, `log level (1-4)
  1 - DEBUG
  2 - INFO (default)
  3 - WARN
  4 - ERROR`)
}

func Execute() {
	rootCmd.AddCommand(get.NewCommand())
	rootCmd.AddCommand(set.NewCommand())
	rootCmd.AddCommand(reload.NewCommand())
	rootCmd.AddCommand(version.NewCommand())

	carapace.Gen(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		pterm.Printfln("%v", err)
		os.Exit(1)
	}
}
