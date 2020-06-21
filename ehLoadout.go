package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.LoadoutEvent.String()] = ehLoadout
}

func ehLoadout(ext *Extension, e events.Event) {
	evt := e.(*journal.Loadout)
	ext.EdState.Write(func() error {
		cmdr := ext.EdState.Cmdr
		ship := cmdr.SetShip(evt.ShipID)
		ship.Type = evt.Ship
		ship.Ident = evt.ShipIdent
		ship.Name = evt.ShipName
		return nil
	})
}