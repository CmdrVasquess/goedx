package edgx

import (
	"time"

	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
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
	ext.EdState.Write(func() error {
		sellShip(ext.EdState.MustCommander(), evt.Time, &evt.ShipSale)
		return nil
	})
	return 0
}
