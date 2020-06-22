package journal

import "github.com/CmdrVasquess/edgx/events"

type undockedT string

const UndockedEvent = undockedT("Undocked")

func (t undockedT) New() events.Event { return new(Undocked) }
func (t undockedT) String() string    { return string(t) }

type Undocked struct {
	events.Common
	StationName string
	StationType string
}

func init() {
	events.RegisterType(string(UndockedEvent), UndockedEvent)
}
