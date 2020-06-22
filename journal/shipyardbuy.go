package journal

import "github.com/CmdrVasquess/edgx/events"

type shipyardbuyT string

const ShipyardBuyEvent = shipyardbuyT("ShipyardBuy")

func (t shipyardbuyT) New() events.Event { return new(ShipyardBuy) }
func (t shipyardbuyT) String() string    { return string(t) }

type ShipyardBuy struct {
	events.Common
	StoreShipID int
	ShipPrice   int64
}

func init() {
	events.RegisterType(string(ShipyardBuyEvent), ShipyardBuyEvent)
}
