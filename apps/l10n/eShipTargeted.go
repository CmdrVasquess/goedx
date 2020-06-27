package l10n

import (
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/journal"
)

func (loc *Locales) finishShipTargeted(evt *journal.ShipTargeted, chg goedx.Change) {
	if evt.ShipL7d != "" {
		loc.shiptype[evt.Ship] = evt.ShipL7d
	}
}
