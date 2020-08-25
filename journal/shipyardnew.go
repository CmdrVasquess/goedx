package journal

import "github.com/CmdrVasquess/goedx/events"

type shipyardnewT string

const ShipyardNewEvent = shipyardnewT("ShipyardNew")

func (t shipyardnewT) New() events.Event { return new(ShipyardNew) }
func (t shipyardnewT) String() string    { return string(t) }

type ShipyardNew struct {
	events.Common
	NewShipID int
	ShipType  string
}

func (_ *ShipyardNew) EventType() events.Type { return ShipyardNewEvent }

func init() {
	events.RegisterType(string(ShipyardNewEvent), ShipyardNewEvent)
}
