package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.SetUserShipNameEvent.String()] = ehSetUserShipName
}

func ehSetUserShipName(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.SetUserShipName)
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		ship := cmdr.GetShip(evt.ShipID)
		ship.Ident = evt.UserShipId
		ship.Name = evt.UserShipName
		return nil
	}))
	return chg
}
