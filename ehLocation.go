package goedx

import (
	"github.com/CmdrVasquess/goedx/att"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	evtHdlrs[journal.LocationEvent.String()] = ehLocation
}

func ehLocation(ed *EDState, e events.Event) (chg att.Change, err error) {
	evt := e.(*journal.Location)
	sys := ed.Galaxy.EdgxSystem(
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
			Type:   evt.StationType,
			Docked: evt.Docked,
		}
	default:
		loc = sys
	}
	err = ed.WrLocked(func() error {
		ed.Loc.Location = loc
		chg = ChgLocation
		return nil
	})
	return chg, err
}
