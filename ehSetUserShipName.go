package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.SetUserShipNameEvent.String()] = ehSetUserShipName
}

func ehSetUserShipName(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.SetUserShipName)
	ext.EdState.Write(func() error {
		cmdr := ext.EdState.MustCommander()
		ship := cmdr.GetShip(evt.ShipID)
		ship.Ident = evt.UserShipId
		ship.Name = evt.UserShipName
		return nil
	})
	return chg
}
