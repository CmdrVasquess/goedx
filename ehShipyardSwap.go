package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShipyardSwapEvent.String()] = ehShipyardSwap
}

func ehShipyardSwap(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.ShipyardSwap)
	Must(ext.EdState.WriteCmdr(func(cmdr *Commander) error {
		if cmdr.ShipID != evt.StoreShipID {
			log.Warna("current `ship` differs from `ship to be stored`",
				cmdr.ShipID,
				evt.StoreShipID)
		}
		cmdr.StoreCurrentShip(evt.StoreShipID)
		cmdr.SetShip(evt.ShipID)
		return nil
	}))
	return 0
}
