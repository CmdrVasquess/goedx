package edgx

import (
	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/edgx/journal"
)

func init() {
	stdEvtHdlrs[journal.CommanderEvent.String()] = ehCommander
}

func ehCommander(ext *Extension, e events.Event) {
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
			cmdr = new(Commander)
			ext.EdState.Cmdr = cmdr
		}
		cmdr.FID = evt.FID
		cmdr.Name = evt.Name
		return nil
	})
}
