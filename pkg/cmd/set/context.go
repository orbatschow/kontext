package set

import (
	"os"

	"github.com/orbatschow/kontext/pkg/completion"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

func runSetContext(_ *cobra.Command, args []string) {
	log := logger.New()
	var contextName string
	if len(args) == 0 {
		contextName = ""
	} else {
		contextName = args[0]
	}

	client, err := context.New()
	if err != nil {
		log.Error("failed to create context client", "err", err)
		os.Exit(1)
	}

	err = client.Set(contextName)
	if err != nil {
		log.Error("failed to set context", "err", err)
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

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context [name]",
		Short: "Switch to a different Kubernetes context",
		Long: `Change the active Kubernetes context within the current group.

Without a name, opens an interactive selection dialog showing contexts from the active group.
With a name, immediately switches to that context.

Use '-' as the context name to switch back to the previously active context.`,
		PreRun: init_,
		Run:    runSetContext,
	}

	carapace.Gen(cmd).PositionalAnyCompletion(completion.Context())

	return cmd
}
