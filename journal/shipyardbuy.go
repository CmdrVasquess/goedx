package journal

import "github.com/CmdrVasquess/goedx/events"

type shipyardbuyT string

const ShipyardBuyEvent = shipyardbuyT("ShipyardBuy")

func (t shipyardbuyT) New() events.Event { return new(ShipyardBuy) }
func (t shipyardbuyT) String() string    { return string(t) }

type ShipyardBuy struct {
	events.Common
	ShipPrice   int64
	StoreShipID int
	SellShipID  int
}

func (_ *ShipyardBuy) EventType() events.Type { return ShipyardBuyEvent }

func init() {
	events.RegisterType(string(ShipyardBuyEvent), ShipyardBuyEvent)
}
