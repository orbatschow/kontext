package logger

import "github.com/pterm/pterm"

func New() *pterm.Logger {
	return pterm.DefaultLogger.WithLevel(pterm.LogLevelDebug)
}
