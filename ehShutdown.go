package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ShutdownEvent.String()] = ehShutdown
}

func ehShutdown(ext *Extension, e events.Event) (chg Change) {
	Must(ext.EDState.Write(func() error {
		if ext.ShutdownLogsOut {
			ext.SwitchCommander("", "")
		} else if ext.EDState.Cmdr != nil && ext.CmdrFile != nil {
			cmdrFile := ext.CmdrFile(ext.EDState.Cmdr)
			if err := ext.EDState.Cmdr.Save(cmdrFile); err != nil {
				return err
			}
		}
		return nil
	}))
	return 0
}
