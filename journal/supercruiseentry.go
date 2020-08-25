package journal

import "github.com/CmdrVasquess/goedx/events"

type supercruiseentryT string

const SupercruiseEntryEvent = supercruiseentryT("SupercruiseEntry")

func (t supercruiseentryT) New() events.Event { return new(SupercruiseEntry) }
func (t supercruiseentryT) String() string    { return string(t) }

type SupercruiseEntry struct {
	events.Common
	StarSystem    string
	SystemAddress uint64
}

func (_ *SupercruiseEntry) EventType() events.Type { return SupercruiseEntryEvent }

func init() {
	events.RegisterType(string(SupercruiseEntryEvent), SupercruiseEntryEvent)
}
