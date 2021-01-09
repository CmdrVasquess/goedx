package goedx

import (
	"github.com/CmdrVasquess/goedx/att"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	evtHdlrs[journal.FileheaderEvent.String()] = ehFileheader
}

func ehFileheader(ed *EDState, e events.Event) (chg att.Change, err error) {
	evt := e.(*journal.Fileheader)
	err = ed.WrLocked(func() error {
		ed.SetEDVersion(evt.GameVersion)
		ed.SetLanguage(evt.Language)
		return nil
	})
	return chg, err
}
