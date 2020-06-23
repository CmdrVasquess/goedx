package journal

import "github.com/CmdrVasquess/goedx/events"

type reputationT string

const ReputationEvent = reputationT("Reputation")

func (t reputationT) New() events.Event { return new(Reputation) }
func (t reputationT) String() string    { return string(t) }

type Reputation struct {
	events.Common
	Alliance    int
	Empire      int
	Federation  int
	Independent int
}

func init() {
	events.RegisterType(string(ReputationEvent), ReputationEvent)
}
