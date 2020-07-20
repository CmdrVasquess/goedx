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
		chg |= ship.Ident.Set(evt.UserShipId, ChgShip)
		chg |= ship.Name.Set(evt.UserShipName, ChgShip)
		return nil
	}))
	return chg
}
