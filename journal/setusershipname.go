package journal

import "github.com/CmdrVasquess/goedx/events"

type setusershipnameT string

const SetUserShipNameEvent = setusershipnameT("SetUserShipName")

func (t setusershipnameT) New() events.Event { return new(SetUserShipName) }
func (t setusershipnameT) String() string    { return string(t) }

type SetUserShipName struct {
	events.Common
	Ship         string
	ShipID       int
	UserShipId   string
	UserShipName string
}

func init() {
	events.RegisterType(string(SetUserShipNameEvent), SetUserShipNameEvent)
}
