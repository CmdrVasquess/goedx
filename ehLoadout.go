package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.LoadoutEvent.String()] = ehLoadout
}

func ehLoadout(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Loadout)
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		ship, switched := cmdr.SetShip(evt.ShipID)
		if switched {
			chg |= ChgShip
		}
		ship.Type = evt.Ship
		chg |= ship.Cargo.Set(evt.CargoCapacity, ChgShip)
		chg |= ship.Ident.Set(evt.ShipIdent, ChgShip)
		chg |= ship.Name.Set(evt.ShipName, ChgShip)
		chg |= ship.MaxRange.Set(evt.MaxJumpRange, ChgShip)
		if ship.MaxRange < ship.MaxJump {
			chg |= ship.MaxJump.Set(0, ChgShip)
		}
		return nil
	}))
	return chg
}
