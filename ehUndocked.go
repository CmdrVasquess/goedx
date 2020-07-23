package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.UndockedEvent.String()] = ehUndocked
}

func ehUndocked(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Undocked)
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		if port := cmdr.At.Port(); port == nil {
			port := &Port{
				Name:   evt.StationName,
				Type:   evt.StationType,
				Docked: false,
			}
			cmdr.At.Location = port
		} else {
			port.Docked = false
			if port.Name != evt.StationName {
				port.Name = evt.StationName
				port.Type = evt.StationType
				port.Sys = nil
			}
		}
		chg = ChgLocation
		return nil
	}))
	return chg
}
