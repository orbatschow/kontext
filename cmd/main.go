package main

import (
	"log"

	"github.com/orbatschow/kontext/pkg/backup"
	"github.com/orbatschow/kontext/pkg/cmd"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
)

func main() {
	// load config
	err := config.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// initialize logger
	logger.Init(config.Get())

	// initialize state
	err = state.Init(config.Get())
	if err != nil {
		log.Fatal(err.Error())
	}

	// create backup
	err = backup.Create(config.Get())
	if err != nil {
		log.Fatal(err.Error())
	}

	// run kontext
	cmd.Execute()
}
