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
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		ship := cmdr.SetShip(evt.ShipID)
		ship.Type = evt.Ship
		ship.Ident = evt.ShipIdent
		ship.Name = evt.ShipName
		return nil
	}))
	return chg
}
