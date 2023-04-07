package set

import (
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func newSetGroupCommand(_ *cobra.Command, args []string) {
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
}

func NewSetContextCommand(_ *cobra.Command, args []string) {
	log := logger.New()
	var contextName string
	if len(args) == 0 {
		contextName = ""
	} else {
		contextName = args[0]
	}

	client, err := context.New()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err = client.Set(contextName)
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

	setGroupCommand := &cobra.Command{
		Use:   "group [name]",
		Short: "set group to active",
		Run:   newSetGroupCommand,
	}

	setContextCommand := &cobra.Command{
		Use:   "context [name]",
		Short: "set context to active",
		Run:   NewSetContextCommand,
	}

	cmd.AddCommand(setGroupCommand)
	cmd.AddCommand(setContextCommand)

	return cmd
}
