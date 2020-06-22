package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.DockedEvent.String()] = ehDocked
}

func ehDocked(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Docked)
	sys, _ := ext.Galaxy.EdgxSystem(evt.SystemAddress, evt.StarSystem, nil)
	loc := &Port{
		Sys:    sys,
		Name:   evt.StationName,
		Docked: true,
	}
	ext.EdState.Write(func() error {
		cmdr := ext.EdState.MustCommander()
		cmdr.Loc.Location = loc
		return nil
	})
	return chg
}
