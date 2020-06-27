package journal

import "github.com/CmdrVasquess/goedx/events"

type shiptargetedT string

const ShipTargetedEvent = shiptargetedT("ShipTargeted")

func (t shiptargetedT) New() events.Event { return new(ShipTargeted) }
func (t shiptargetedT) String() string    { return string(t) }

type ShipTargeted struct {
	events.Common
	Ship    string
	ShipL7d string `json:"Ship_Localised"`
}

func init() {
	events.RegisterType(string(ShipTargetedEvent), ShipTargetedEvent)
}
