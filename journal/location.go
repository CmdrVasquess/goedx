package journal

import "github.com/CmdrVasquess/edgx/events"

type locationT string

const LocationEvent = locationT("Location")

func (t locationT) New() events.Event { return new(Commander) }
func (t locationT) String() string    { return string(t) }

type Location struct {
	events.Common
	StarSystem    string
	SystemAddress uint64
	StarPos       [3]float64
	Docked        bool
	StationName   string
	StationType   string
}

func init() {
	events.RegisterType(string(LocationEvent), LocationEvent)
}
