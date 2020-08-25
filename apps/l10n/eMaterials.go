package l10n

import (
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/journal"
)

func (loc *Locales) finishMaterials(evt *journal.Materials, chg goedx.Change) {
	for _, mat := range evt.Raw {
		if mat.NameL7d != "" {
			loc.MatsRaw[normKey(mat.Name)] = mat.NameL7d
		}
	}
	for _, mat := range evt.Manufactured {
		if mat.NameL7d != "" {
			loc.MatsMan[normKey(mat.Name)] = mat.NameL7d
		}
	}
	for _, mat := range evt.Encoded {
		if mat.NameL7d != "" {
			loc.MatsEnc[normKey(mat.Name)] = mat.NameL7d
		}
	}
}
