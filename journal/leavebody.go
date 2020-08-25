package journal

import "github.com/CmdrVasquess/goedx/events"

type leavebodyT string

const LeaveBodyEvent = leavebodyT("LeaveBody")

func (t leavebodyT) New() events.Event { return new(LeaveBody) }
func (t leavebodyT) String() string    { return string(t) }

type LeaveBody struct {
	events.Common
	StarSystem    string
	SystemAddress uint64
	Body          string
	BodyID        int
}

func (_ *LeaveBody) EventType() events.Type { return LeaveBodyEvent }

func init() {
	events.RegisterType(string(LeaveBodyEvent), LeaveBodyEvent)
}
