package goedx

import (
	"time"

	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShipyardSellEvent.String()] = ehShipyardSell
}

func sellShip(cmdr *Commander, t time.Time, evt *journal.ShipSale) {
	ship := cmdr.FindShip(evt.SellShipID)
	if ship == nil {
		log.Warna("mssing ship `id` for `event`",
			evt.SellShipID,
			journal.SellShipOnRebuyEvent)
		return
	}
	ship.Sold = new(time.Time)
	*ship.Sold = t
	if cmdr.ShipID == evt.SellShipID {
		cmdr.ShipID = -1
		cmdr.inShip = nil
	}
}

func ehShipyardSell(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.ShipyardSell)
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		sellShip(cmdr, evt.Time, &evt.ShipSale)
		return nil
	}))
	return 0
}
