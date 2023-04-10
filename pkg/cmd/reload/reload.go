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

			file, err := os.OpenFile(config.Get().Global.Kubeconfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
			err = kubeconfig.Write(file, client.APIConfig)
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
