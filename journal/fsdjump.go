package journal

import "github.com/CmdrVasquess/edgx/events"

type fsdjumpT string

const FSDJumpEvent = fsdjumpT("FSDJump")

func (t fsdjumpT) New() events.Event { return new(FSDJump) }
func (t fsdjumpT) String() string    { return string(t) }

type FSDJump struct {
	events.Common
	SystemAddress uint64
	StarSystem    string
	StarPos       [3]float32
	Population    int64
	Body          string
	BodyID        int
}

func init() {
	events.RegisterType(string(FSDJumpEvent), FSDJumpEvent)
}
