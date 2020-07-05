package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.FSDJumpEvent.String()] = ehFSDJump
}

func ehFSDJump(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.FSDJump)
	sys, _ := ext.Galaxy.EdgxSystem(
		evt.SystemAddress,
		evt.StarSystem,
		evt.StarPos[:],
		evt.Time,
	)
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		cmdr.Jump(evt.SystemAddress, evt.Time)
		cmdr.At.Location = sys
		// TODO be more precise
		return nil
	}))
	return chg
}
