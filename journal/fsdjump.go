package journal

import "github.com/CmdrVasquess/goedx/events"

type fsdjumpT string

const FSDJumpEvent = fsdjumpT("FSDJump")

func (t fsdjumpT) New() events.Event { return new(FSDJump) }
func (t fsdjumpT) String() string    { return string(t) }

type FSDJump struct {
	events.Common
	SystemAddress          uint64
	StarSystem             string
	StarPos                [3]float32
	JumpDist               float32
	Population             int64
	Body                   string
	BodyID                 int
	SystemEconomy          string
	SystemEconomyL7d       string `json:"SystemEconomy_Localised"`
	SystemSecondEconomy    string
	SystemSecondEconomyL7d string `json:"SystemSecondEconomy_Localised"`
	SystemSecurity         string
	SystemSecurityL7d      string `json:"SystemSecurity_Localised"`
}

func init() {
	events.RegisterType(string(FSDJumpEvent), FSDJumpEvent)
}
