package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShutdownEvent.String()] = ehShutdown
}

func ehShutdown(ext *Extension, e events.Event) (chg Change) {
	Must(ext.EdState.Write(func() error {
		ext.SwitchCommander("", "")
		return nil
	}))
	return 0
}
