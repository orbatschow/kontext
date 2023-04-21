package set

import (
	"log"
	"os"

	"github.com/orbatschow/kontext/pkg/backup"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/kubeconfig"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func Init(_ *cobra.Command, _ []string) {
	// load currentConfig
	configClient := &config.Client{
		File: config.DefaultConfigPath,
	}
	currentConfig, err := configClient.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize logger
	logger.Init(currentConfig)

	// initialize currentState
	err = state.Init(currentConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// read the current currentState
	currentState, err := state.Read(currentConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create backup
	backupReconciler := backup.Reconciler{
		Config: currentConfig,
		State:  currentState,
	}
	err = backupReconciler.Reconcile()
	if err != nil {
		log.Fatal(err.Error())
	}
}

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

	file, err := os.OpenFile(client.Config.Global.Kubeconfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err = kubeconfig.Write(file, client.APIConfig)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err = state.Write(client.Config, client.State)
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

	file, err := os.OpenFile(client.Config.Global.Kubeconfig, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	err = kubeconfig.Write(file, client.APIConfig)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err = state.Write(client.Config, client.State)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [name]",
		Short: "set [context|group] [name]",
		Long: `Invoking this command without a parameter will spawn an interactive selection dialog.
When providing a context name, the switch will be performed immediately.
'-' is a reserved context name, that will cause a switch to the previously active context.
If neither, context nor group is specified, this command will set the context.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(1)
		},
	}

	setGroupCommand := &cobra.Command{
		Use:   "group [name]",
		Short: "set group to active",
		Long: `Invoking this command without a parameter will spawn an interactive selection dialog. 
When providing a group name, the switch will be performed immediately.
'-' is a reserved group name, that will cause a switch to the previously active group.
		`,
		PreRun: Init,
		Run:    newSetGroupCommand,
	}

	setContextCommand := &cobra.Command{
		Use:   "context [name]",
		Short: "set context to active",
		Long: `Invoking this command without a parameter will spawn an interactive selection dialog.
When providing a context name, the switch will be performed immediately.
'-' is a reserved context name, that will cause a switch to the previously active context.
		`,
		PreRun: Init,
		Run:    NewSetContextCommand,
	}

	cmd.AddCommand(setGroupCommand)
	cmd.AddCommand(setContextCommand)

	return cmd
}
