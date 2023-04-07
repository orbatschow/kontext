package reload

import (
	"os"

	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "reload the active group",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()

			client, err := group.New()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}

			err = client.Reload()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}

			err = state.Write(client.State)
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		},
	}
	return cmd
}
