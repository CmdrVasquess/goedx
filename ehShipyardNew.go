package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShipyardNewEvent.String()] = ehShipyardNew
}

func ehShipyardNew(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.ShipyardNew)
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		cmdr.StoreCurrentShip(0)
		ship := cmdr.GetShip(evt.NewShipID)
		ship.Type = evt.ShipType
		cmdr.inShip = ship
		return nil
	}))
	return 0
}
