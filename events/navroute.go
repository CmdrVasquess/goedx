package events

import (
	"github.com/CmdrVasquess/watched"
)

type navrouteT string

const NavRouteEvent = statusT(watched.StatNavRouteName)

func (t navrouteT) New() Event     { return new(NavRoute) }
func (t navrouteT) String() string { return string(t) }

type WayPoint struct {
	StarSystem    string
	SystemAddress uint64
	StarPos       [3]float32
	StarClass     string
}

type NavRoute struct {
	Common
	Route []WayPoint
}

func (_ *NavRoute) EventType() Type { return NavRouteEvent }

func init() {
	RegisterType(string(StatusEvent), StatusEvent)
}
