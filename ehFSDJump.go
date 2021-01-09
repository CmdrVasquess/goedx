package goedx

import (
	"github.com/CmdrVasquess/goedx/att"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	evtHdlrs[journal.FSDJumpEvent.String()] = ehFSDJump
}

func ehFSDJump(ed *EDState, e events.Event) (chg att.Change, err error) {
	evt := e.(*journal.FSDJump)
	chg = ChgSystem
	sys := ed.Galaxy.EdgxSystem(
		evt.SystemAddress,
		evt.StarSystem,
		evt.StarPos[:],
		evt.Time,
	)
	err = ed.WrLocked(func() error {
		ed.JumpHist.Jump(evt.SystemAddress, evt.Time)
		ed.Loc.Location = sys
		cmdr := ed.Cmdr
		if cmdr.inShip != nil && evt.JumpDist > float32(cmdr.inShip.MaxJump) {
			chg |= cmdr.inShip.MaxJump.Set(evt.JumpDist, ChgShip)
		}
		// TODO be more precise
		return nil
	})
	return chg, err
}
