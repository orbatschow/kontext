package main

import (
	"github.com/orbatschow/kontext/pkg/backup"
	"github.com/orbatschow/kontext/pkg/cmd"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
)

func main() {
	// initialize logger
	log := logger.New()

	// initialize state
	err := state.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// load config
	err = config.Read()
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
