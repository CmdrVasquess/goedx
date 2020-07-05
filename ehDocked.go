package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.DockedEvent.String()] = ehDocked
}

func ehDocked(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Docked)
	sys, _ := ext.Galaxy.EdgxSystem(evt.SystemAddress, evt.StarSystem, nil, evt.Time)
	loc := &Port{
		Sys:    sys,
		Name:   evt.StationName,
		Docked: true,
	}
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		cmdr.At.Location = loc
		chg = ChgLocation
		return nil
	}))
	return chg
}
