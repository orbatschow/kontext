package get

import (
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/orbatschow/kontext/pkg/completion"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/spf13/cobra"
)

func newGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "group [name]",
		Short:  "List or display configured groups",
		Long: `List all configured groups with their source configurations.

If a group name is provided, display detailed information about that specific group including
its sources and context configuration.`,
		PreRun: init_,
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			var groupName string

			if len(args) > 0 {
				groupName = args[0]
			}

			groupClient, err := group.New()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}

			var groups []config.GroupItem
			if len(groupName) != 0 {
				match, err := groupClient.Get(groupName)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				groups = append(groups, *match)
			} else {
				groups = groupClient.Config.Groups.Items
			}

			printer := groupClient.BuildTablePrinter(groups...)
			err = printer.Render()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		},
	}

	carapace.Gen(cmd).PositionalAnyCompletion(completion.Group())

	return cmd
}
