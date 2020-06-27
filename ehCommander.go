package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.CommanderEvent.String()] = ehCommander
}

func ehCommander(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Commander)
	Must(ext.EdState.Write(func() error {
		cmdr := ext.EdState.Cmdr
		if cmdr != nil && cmdr.FID != "" {
			if ext.CmdrFile != nil {
				f := ext.CmdrFile(cmdr)
				if err := cmdr.Save(f); err != nil {
					log.Errore(err)
				}
			}
		}
		cmdr = NewCommander(evt.FID)
		if evt.FID != "" && ext.CmdrFile != nil {
			f := ext.CmdrFile(cmdr)
			if err := cmdr.Load(f); err != nil {
				log.Errore(err)
			}
		}
		cmdr.FID = evt.FID
		cmdr.Name = evt.Name
		ext.EdState.Cmdr = cmdr
		return nil
	}))
	return chg
}
