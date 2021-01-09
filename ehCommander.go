package goedx

import (
	"github.com/CmdrVasquess/goedx/att"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	evtHdlrs[journal.CommanderEvent.String()] = ehCommander
}

func ehCommander(ed *EDState, e events.Event) (chg att.Change, err error) {
	evt := e.(*journal.Commander)
	err = ed.WrLocked(func() error {
		ed.SwitchCommander(evt.FID, evt.Name)
		return nil
	})
	return chg, err
}
