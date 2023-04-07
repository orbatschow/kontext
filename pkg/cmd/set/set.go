package set

import (
	"os"

	kubectx "github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/spf13/cobra"
)

func newSetGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group [name]",
		Short: "set group to active",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			err := group.Set(args[0])
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}
		},
	}
	return cmd
}

func newSetContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context [name]",
		Short: "set context to active",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			err := kubectx.Set(args[0])
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
		Use:   "set",
		Short: "set [context|group] [name]",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(newSetGroupCommand())
	cmd.AddCommand(newSetContextCommand())

	return cmd
}
