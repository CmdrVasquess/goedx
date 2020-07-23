package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.SupercruiseEntryEvent.String()] = ehSupercruiseEntry
}

func ehSupercruiseEntry(ext *Extension, _ events.Event) (chg Change) {
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		if cmdr.At.Location == nil {
			return nil
		}
		cmdr.At.Location = cmdr.At.Location.System()
		chg = ChgLocation
		return nil
	}))
	return chg
}
