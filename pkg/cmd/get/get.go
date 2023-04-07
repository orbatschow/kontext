package get

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	kubectx "github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/spf13/cobra"
)

func newGetGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group [name]",
		Short: "get groups, optionally filtered by name",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			kontextConfig := config.Get()
			var groupName string
			if len(args) > 0 {
				groupName = args[0]
			}
			err := group.Get(cmd, kontextConfig, groupName)
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		},
	}
	return cmd
}

func newGetContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context [name]",
		Short: "get contexts, optionally filtered by name",
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			var contextName string
			if len(args) > 0 {
				contextName = args[0]
			}
			err := kubectx.Get(cmd, contextName)
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
		Use:   "get",
		Short: "get [context|group] [name], defaults to context",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(newGetGroupCommand())
	cmd.AddCommand(newGetContextCommand())
	return cmd
}
