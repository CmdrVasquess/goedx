package journal

import "github.com/CmdrVasquess/edgx/events"

type sellshiponrebuyT string

const SellShipOnRebuyEvent = sellshiponrebuyT("SellShipOnRebuy")

func (t sellshiponrebuyT) New() events.Event { return new(SellShipOnRebuy) }
func (t sellshiponrebuyT) String() string    { return string(t) }

type SellShipOnRebuy struct {
	events.Common
	ShipSale
}

func init() {
	events.RegisterType(string(SellShipOnRebuyEvent), SellShipOnRebuyEvent)
}
