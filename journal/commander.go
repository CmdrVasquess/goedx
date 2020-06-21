package journal

import "github.com/CmdrVasquess/edgx/events"

type commanderT string

const CommanderEvent = commanderT("Commander")

func (t commanderT) New() events.Event { return new(Commander) }
func (t commanderT) String() string    { return string(t) }

type Commander struct {
	events.Common
	FID  string
	Name string
}

func init() {
	events.RegisterType(string(CommanderEvent), CommanderEvent)
}
