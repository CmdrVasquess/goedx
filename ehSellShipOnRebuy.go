package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.SellShipOnRebuyEvent.String()] = ehSellShipOnRebuy
}

func ehSellShipOnRebuy(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.SellShipOnRebuy)
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		sellShip(cmdr, evt.Time, &evt.ShipSale)
		return nil
	}))
	return 0
}
