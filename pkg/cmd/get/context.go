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
		Short:  "get contexts, optionally filtered by name",
		PreRun: init_,
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			var contextName string
			if len(args) > 0 {
				contextName = args[0]
			}

			contextClient, err := context.New()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}

			if len(contextName) != 0 {
				match, err := contextClient.Get(contextName)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				printer := contextClient.BuildTablePrinter(match)
				err = printer.Render()
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
			} else {
				match := contextClient.List()
				printer := contextClient.BuildTablePrinter(match)
				err = printer.Render()
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
			}
		},
	}

	carapace.Gen(cmd).PositionalAnyCompletion(completion.Context())

	return cmd
}
