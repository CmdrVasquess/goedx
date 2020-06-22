package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.CommanderEvent.String()] = ehCommander
}

func ehCommander(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Commander)
	ext.EdState.Write(func() error {
		cmdr := ext.EdState.Cmdr
		if cmdr != nil {
			if ext.CmdrFile != nil {
				f := ext.CmdrFile(cmdr)
				if err := cmdr.Save(f); err != nil {
					log.Errore(err)
				}
			}
		}
		cmdr = new(Commander)
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
	})
	return chg
}
