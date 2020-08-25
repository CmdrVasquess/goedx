package journal

import "github.com/CmdrVasquess/goedx/events"

type shipyardswapT string

const ShipyardSwapEvent = shipyardswapT("ShipyardSwap")

func (t shipyardswapT) New() events.Event { return new(ShipyardSwap) }
func (t shipyardswapT) String() string    { return string(t) }

type ShipyardSwap struct {
	events.Common
	ShipID      int
	StoreShipID int
}

func (_ *ShipyardSwap) EventType() events.Type { return ShipyardSwapEvent }

func init() {
	events.RegisterType(string(ShipyardSwapEvent), ShipyardSwapEvent)
}
