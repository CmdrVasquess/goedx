package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.LocationEvent.String()] = ehLocation
}

func ehLocation(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Location)
	sys, _ := ext.Galaxy.EdgxSystem(
		evt.SystemAddress,
		evt.StarSystem,
		evt.StarPos[:],
		evt.Time,
	)
	var loc Location
	switch {
	case evt.StationName != "":
		loc = &Port{
			Sys:    sys,
			Name:   evt.StationName,
			Docked: evt.Docked,
		}
	default:
		loc = sys
	}
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		cmdr.At.Location = loc
		chg = ChgLocation
		return nil
	}))
	return chg
}
