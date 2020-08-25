package journal

import "github.com/CmdrVasquess/goedx/events"

type dockedT string

const DockedEvent = dockedT("Docked")

func (t dockedT) New() events.Event { return new(Docked) }
func (t dockedT) String() string    { return string(t) }

type Docked struct {
	events.Common
	SystemAddress  uint64
	StarSystem     string
	StationName    string
	StationType    string
	DistFromStarLS float64
}

func (_ *Docked) EventType() events.Type { return DockedEvent }

func init() {
	events.RegisterType(string(DockedEvent), DockedEvent)
}
