package l10n

import (
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/journal"
)

func (loc *Locales) finishFSDJump(evt *journal.FSDJump, chg goedx.Change) {
	if evt.SystemSecondEconomyL7d != "" {
		key := normKey(evt.SystemSecondEconomy)
		loc.Economies[key] = evt.SystemSecondEconomyL7d
	}
	if evt.SystemEconomyL7d != "" {
		key := normKey(evt.SystemEconomy)
		loc.Economies[key] = evt.SystemEconomyL7d
	}
	if evt.SystemSecurityL7d != "" {
		key := normKey(evt.SystemSecurity)
		loc.Securities[key] = evt.SystemSecurityL7d
	}
}
