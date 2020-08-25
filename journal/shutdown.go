package journal

import "github.com/CmdrVasquess/goedx/events"

type shutdownT string

const ShutdownEvent = shutdownT("Shutdown")

func (t shutdownT) New() events.Event { return new(Shutdown) }
func (t shutdownT) String() string    { return string(t) }

type Shutdown struct{ events.Common }

func (_ *Shutdown) EventType() events.Type { return ShutdownEvent }

func init() {
	events.RegisterType(string(ShutdownEvent), ShutdownEvent)
}
