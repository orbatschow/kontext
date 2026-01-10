package get

import (
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/orbatschow/kontext/pkg/completion"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/spf13/cobra"
)

func newContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "context [name]",
		Short:  "List contexts from the active group",
		Long: `List all Kubernetes contexts available in the currently active group.

If a context name is provided, display detailed information about that specific context.`,
		PreRun: init_,
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			var contextName string
			if len(args) > 0 {
				contextName = args[0]
			}

			contextClient, err := context.New()
			if err != nil {
				log.Error("failed to create context client", "err", err)
				os.Exit(1)
			}

			if len(contextName) != 0 {
				match, err := contextClient.Get(contextName)
				if err != nil {
					log.Error("failed to get context", "err", err)
					os.Exit(1)
				}
				printer := contextClient.BuildTablePrinter(match)
				err = printer.Render()
				if err != nil {
					log.Error("failed to render contexts", "err", err)
					os.Exit(1)
				}
			} else {
				match := contextClient.List()
				printer := contextClient.BuildTablePrinter(match)
				err = printer.Render()
				if err != nil {
					log.Error("failed to render contexts", "err", err)
					os.Exit(1)
				}
			}
		},
	}

	carapace.Gen(cmd).PositionalAnyCompletion(completion.Context())

	return cmd
}
