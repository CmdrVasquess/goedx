package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.SupercruiseExitEvent.String()] = ehSupercruiseExit
}

func ehSupercruiseExit(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.SupercruiseExit)
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		if cmdr.At.Location == nil {
			return nil
		}
		if evt.BodyType != "Station" {
			return nil
		}
		sys := cmdr.At.Location.System()
		port := &Port{
			Sys:  sys,
			Name: evt.Body,
		}
		cmdr.At.Location = port
		chg = ChgLocation
		return nil
	}))
	return chg
}
