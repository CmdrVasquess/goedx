package journal

import "github.com/CmdrVasquess/goedx/events"

type fssdiscoveryscanT string

const FSSDiscoveryScanEvent = fssdiscoveryscanT("FSSDiscoveryScan")

func (t fssdiscoveryscanT) New() events.Event { return new(FSSDiscoveryScan) }
func (t fssdiscoveryscanT) String() string    { return string(t) }

type FSSDiscoveryScan struct {
	events.Common
	SystemAddress uint64
	SystemName    string
	BodyCount     int
	NonBodyCount  int
}

func (_ *FSSDiscoveryScan) EventType() events.Type { return FSSDiscoveryScanEvent }

func init() {
	events.RegisterType(string(FSSDiscoveryScanEvent), FSSDiscoveryScanEvent)
}
