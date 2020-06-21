package journal

import "github.com/CmdrVasquess/edgx/events"

type loadoutT string

const LoadoutEvent = loadoutT("Loadout")

func (t loadoutT) New() events.Event { return new(Loadout) }
func (t loadoutT) String() string    { return string(t) }

type Loadout struct {
	events.Common
	Ship      string
	ShipID    int
	ShipName  string
	ShipIdent string
}

func init() {
	events.RegisterType(string(LoadoutEvent), LoadoutEvent)
}