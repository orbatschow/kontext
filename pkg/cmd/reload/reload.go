package reload

import (
	"os"

	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func newReloadGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
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

func newReloadContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "reload the active context",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()

			client, err := context.New()
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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "reload <context|group>",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(newReloadGroupCommand())
	cmd.AddCommand(newReloadContextCommand())
	return cmd
}
