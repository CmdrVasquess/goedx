package l10n

import (
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/journal"
)

func (loc *Locales) finishFSDJump(evt *journal.FSDJump, chg goedx.Change) {
	if evt.SystemSecondEconomyL7d != "" {
		loc.economy[evt.SystemSecondEconomy] = evt.SystemSecondEconomyL7d
	}
	if evt.SystemEconomyL7d != "" {
		loc.economy[evt.SystemEconomy] = evt.SystemEconomyL7d
	}
	if evt.SystemSecurityL7d != "" {
		loc.security[evt.SystemSecurity] = evt.SystemSecurityL7d
	}
}
