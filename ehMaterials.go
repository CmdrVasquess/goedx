package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.MaterialsEvent.String()] = ehMaterials
}

func ehMaterials(ext *Extension, e events.Event) (chg Change) {
	cpMats := func(sMats map[string]*Material, jMats []journal.Material) {
		for j := range jMats {
			jm := &jMats[j]
			if sm, have := sMats[jm.Name]; have {
				sm.Stock = jm.Count
			} else {
				sMats[jm.Name] = &Material{Stock: jm.Count}
			}
		}
	}
	evt := e.(*journal.Materials)
	Must(ext.EDState.WriteCmdr(func(cmdr *Commander) error {
		cpMats(cmdr.Mats.Raw, evt.Raw)
		cpMats(cmdr.Mats.Man, evt.Manufactured)
		cpMats(cmdr.Mats.Enc, evt.Encoded)
		return nil
	}))
	return chg
}
