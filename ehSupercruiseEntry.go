package goedx

import (
	"github.com/CmdrVasquess/goedx/att"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	evtHdlrs[journal.SupercruiseEntryEvent.String()] = ehSupercruiseEntry
}

func ehSupercruiseEntry(ed *EDState, _ events.Event) (chg att.Change, err error) {
	err = ed.WrLocked(func() error {
		if ed.Loc.Location == nil {
			return nil
		}
		ed.Loc.Location = ed.Loc.Location.System()
		chg = ChgLocation
		return nil
	})
	return chg, err
}
