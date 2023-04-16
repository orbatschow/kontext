package get

import (
	"log"
	"os"

	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/context"
	"github.com/orbatschow/kontext/pkg/group"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
	"github.com/spf13/cobra"
)

func Init(_ *cobra.Command, _ []string) {
	// load config
	configClient := &config.Client{
		Path: config.DefaultConfigPath,
	}
	currentConfig, err := configClient.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize logger
	logger.Init(currentConfig)

	// initialize state
	err = state.Init(currentConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func newGetGroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "group [name]",
		Short:  "get groups, optionally filtered by name",
		PreRun: Init,
		Run: func(cmd *cobra.Command, args []string) {
			log := logger.New()
			var groupName string

			if len(args) > 0 {
				groupName = args[0]
			}

			groupClient, err := group.New()
			if err != nil {
				log.Error(err.Error())
				os.Exit(1)
			}

			var groups []config.Group
			if len(groupName) != 0 {
				match, err := groupClient.Get(groupName)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				groups = append(groups, *match)
			} else {
				groups = groupClient.Config.Groups
			}

			err = groupClient.Print(groups...)
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
		Use:    "context [name]",
		Short:  "get contexts, optionally filtered by name",
		PreRun: Init,
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

			// if a context name is given, find it and render the result
			if len(contextName) != 0 {
				match, err := contextClient.Get(contextName)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				err = contextClient.Print(match)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
			} else {
				// if no context name is given, find all groups and render the result
				match := contextClient.List()
				err = contextClient.Print(match)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
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
