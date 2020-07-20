package goedx

import (
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

func init() {
	stdEvtHdlrs[journal.ScanEvent.String()] = ehScan
}

func ehScan(ext *Extension, e events.Event) (chg Change) {
	evt := e.(*journal.Scan)
	ext.EDState.Write(func() error {
		cmdr := ext.EDState.Cmdr
		if cmdr != nil {
			for _, mat := range evt.Materials {
				if mstat := cmdr.RawMatStats[mat.Name]; mstat == nil {
					cmdr.RawMatStats[mat.Name] = &RawMatStats{
						Min:   mat.Percent,
						Max:   mat.Percent,
						Sum:   float64(mat.Percent),
						Count: 1,
					}
				} else {
					if mat.Percent < mstat.Min {
						mstat.Min = mat.Percent
					} else if mat.Percent > mstat.Max {
						mstat.Max = mat.Percent
					}
					mstat.Sum += float64(mat.Percent)
					mstat.Count++
				}
			}
		}
		return nil
	})
	return chg
}
