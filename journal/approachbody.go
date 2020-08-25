package journal

import "github.com/CmdrVasquess/goedx/events"

type approachbodyT string

const ApproachBodyEvent = approachbodyT("ApproachBody")

func (t approachbodyT) New() events.Event { return new(ApproachBody) }
func (t approachbodyT) String() string    { return string(t) }

type ApproachBody struct {
	events.Common
	StarSystem    string
	SystemAddress uint64
	Body          string
	BodyID        int
}

func (_ *ApproachBody) EventType() events.Type { return ApproachBodyEvent }

func init() {
	events.RegisterType(string(ApproachBodyEvent), ApproachBodyEvent)
}
