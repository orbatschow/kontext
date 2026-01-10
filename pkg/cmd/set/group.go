package set

import (
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/orbatschow/kontext/pkg/completion"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func runSetGroup(_ *cobra.Command, args []string) {
	log := logger.New()
	var groupName string
	if len(args) == 0 {
		groupName = ""
	} else {
		groupName = args[0]
	}

	client, err := group.New()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err = client.Set(groupName)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	file, err := os.OpenFile(client.Config.Global.Kubeconfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err = kubeconfig.Write(file, client.APIConfig)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err = state.Write(client.Config, client.State)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func newGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group [name]",
		Short: "set group to active",
		Long: `Invoking this command without a parameter will spawn an interactive selection dialog.
When providing a group name, the switch will be performed immediately.
'-' is a reserved group name, that will cause a switch to the previously active group.
		`,
		PreRun: init_,
		Run:    runSetGroup,
	}

	carapace.Gen(cmd).PositionalAnyCompletion(completion.Group())

	return cmd
}
