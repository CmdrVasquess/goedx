package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.SellShipOnRebuyEvent.String()] = ehSellShipOnRebuy
}

func ehSellShipOnRebuy(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.SellShipOnRebuy)
	ext.EdState.Write(func() error {
		sellShip(ext.EdState.MustCommander(), evt.Time, &evt.ShipSale)
		return nil
	})
	return 0
}
