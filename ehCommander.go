package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.CommanderEvent.String()] = ehCommander
}

func ehCommander(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Commander)
	Must(ext.EdState.Write(func() error {
		ext.SwitchCommander(evt.FID, evt.Name)
		return nil
	}))
	return chg
}
