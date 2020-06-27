package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShipTargetedEvent.String()] = ehShipTargeted
}

func ehShipTargeted(ext *Extension, e events.Event) (chg Change) { return 0 }
