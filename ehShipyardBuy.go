package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShipyardBuyEvent.String()] = ehShipyardBuy
}

func ehShipyardBuy(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.ShipyardBuy)
	if evt.SellShipID > 0 {
		Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
			sellShip(cmdr, evt.Time, evt.SellShipID)
			return nil
		}))
	} else if evt.StoreShipID > 0 {
		Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
			cmdr.StoreCurrentShip(evt.StoreShipID)
			return nil
		}))
	}
	return 0
}
