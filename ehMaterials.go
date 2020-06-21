package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.MaterialsEvent.String()] = ehMaterials
}

func ehMaterials(ext *Extension, e events.Event) {
	cpMats := func(sMats map[string]Material, jMats []journal.Material) map[string]Material {
		res := make(map[string]Material)
		for j := range jMats {
			jm := &jMats[j]
			if sm, have := sMats[jm.Name]; have {
				sm.Stock = jm.Count
				res[jm.Name] = sm
			} else {
				res[jm.Name] = Material{Stock: jm.Count}
			}
		}
		return res
	}
	evt := e.(*journal.Materials)
	ext.EdState.Write(func() error {
		cmdr := ext.EdState.Cmdr
		if cmdr == nil {
			panic("materials event without commander")
		}
		cmdr.Mats.Raw = cpMats(cmdr.Mats.Raw, evt.Raw)
		cmdr.Mats.Man = cpMats(cmdr.Mats.Man, evt.Manufactured)
		cmdr.Mats.Enc = cpMats(cmdr.Mats.Enc, evt.Encoded)
		return nil
	})
}