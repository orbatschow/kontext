package logger

import (
	"github.com/pterm/pterm"
)

const DefaultVerbosity = 3

var Verbosity int

func New() *pterm.Logger {
	return pterm.DefaultLogger.WithLevel(pterm.LogLevel(Verbosity))
}
