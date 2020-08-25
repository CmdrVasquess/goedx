package journal

import "github.com/CmdrVasquess/goedx/events"

type locationT string

const LocationEvent = locationT("Location")

func (t locationT) New() events.Event { return new(Location) }
func (t locationT) String() string    { return string(t) }

type Location struct {
	events.Common
	StarSystem    string
	SystemAddress uint64
	StarPos       [3]float32
	Docked        bool
	StationName   string
	StationType   string
}

func (_ *Location) EventType() events.Type { return LocationEvent }

func init() {
	events.RegisterType(string(LocationEvent), LocationEvent)
}
