package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.FileheaderEvent.String()] = ehFileheader
}

func ehFileheader(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Fileheader)
	Must(ext.EDState.Write(func() error {
		ext.EDState.SetEDVersion(evt.GameVersion)
		ext.EDState.SetLanguage(evt.Language)
		return nil
	}))
	return chg
}
