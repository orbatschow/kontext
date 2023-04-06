package main

import (
	"github.com/orbatschow/kontext/pkg/cmd"
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/orbatschow/kontext/pkg/logger"
	"github.com/orbatschow/kontext/pkg/state"
)

func main() {
	// initialize logger
	log := logger.New()

	// initialize state
	err := state.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	// load config
	err = config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	// run kontext
	cmd.Execute()
}
