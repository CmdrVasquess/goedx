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
		if cmdr != nil {
			if ext.CmdrFile != nil {
				f := ext.CmdrFile(cmdr)
				if err := cmdr.Save(f); err != nil {
					log.Errore(err)
				}
			}
		}
		cmdr = NewCommander()
		if ext.CmdrFile != nil {
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
