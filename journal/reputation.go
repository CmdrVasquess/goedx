package journal

import "github.com/CmdrVasquess/goedx/events"

type reputationT string

const ReputationEvent = reputationT("Reputation")

func (t reputationT) New() events.Event { return new(Reputation) }
func (t reputationT) String() string    { return string(t) }

type Reputation struct {
	events.Common
	Alliance    float32
	Empire      float32
	Federation  float32
	Independent float32
}

func (_ *Reputation) EventType() events.Type { return ReputationEvent }

func init() {
	events.RegisterType(string(ReputationEvent), ReputationEvent)
}
