package logger

import (
	"github.com/orbatschow/kontext/pkg/config"
	"github.com/pterm/pterm"
)

var logger *pterm.Logger

func Init(config *config.Config) {
	logger = pterm.DefaultLogger.WithLevel(config.Global.Verbosity)
}

func New() *pterm.Logger {
	if logger == nil {
		return pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	}
	return logger
}
