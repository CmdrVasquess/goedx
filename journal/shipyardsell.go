package journal

import "github.com/CmdrVasquess/goedx/events"

type shipyardsellT string

const ShipyardSellEvent = shipyardsellT("ShipyardSell")

func (t shipyardsellT) New() events.Event { return new(ShipyardSell) }
func (t shipyardsellT) String() string    { return string(t) }

type ShipSale struct {
	SellShipID int
	ShipPrice  int64
}

type ShipyardSell struct {
	events.Common
	ShipSale
}

func init() {
	events.RegisterType(string(ShipyardSellEvent), ShipyardSellEvent)
}
