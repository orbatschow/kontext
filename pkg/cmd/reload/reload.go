package reload

import (
	"os"

	kubectx "github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/spf13/cobra"
)

func newReloadGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "reload the active group",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()

			err := group.Reload(cmd)
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
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()

			err := kubectx.Reload(cmd)
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
