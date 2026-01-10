package reload

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

//revive:disable:var-naming
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

func runReload(_ *cobra.Command, _ []string) {
	log := logger.New()

	client, err := group.New()
	if err != nil {
		log.Error("failed to create group client", "err", err)
		os.Exit(1)
	}

	err = client.Reload()
	if err != nil {
		log.Error("failed to reload group", "err", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(client.Config.Global.Kubeconfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error("failed to open kubeconfig file", "err", err)
		os.Exit(1)
	}
	err = kubeconfig.Write(file, client.APIConfig)
	if err != nil {
		log.Error("failed to write kubeconfig", "err", err)
		os.Exit(1)
	}

	err = state.Write(client.Config, client.State)
	if err != nil {
		log.Error("failed to write state", "err", err)
		os.Exit(1)
	}
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Reload the active group",
		Long: `Re-merge all kubeconfigs from the active group.

This is useful when source kubeconfig files have been updated externally and you
want to refresh the merged kubeconfig without switching groups.`,
		PreRun: init_,
		Run:    runReload,
	}
	return cmd
}
