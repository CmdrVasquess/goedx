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
	Must(ext.EdState.Write(func() error {
		ext.EdState.SetEDVersion(evt.GameVersion)
		ext.EdState.SetLanguage(evt.Language)
		return nil
	}))
	return chg
}
