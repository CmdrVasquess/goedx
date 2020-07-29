package journal

import (
	"git.fractalqb.de/fractalqb/ggja"
	"github.com/CmdrVasquess/goedx/events"
)

type scanT string

const ScanEvent = scanT("Scan")

func (t scanT) New() events.Event { return new(Scan) }
func (t scanT) String() string    { return string(t) }

type ScanMaterial struct {
	Name    string
	Percent float32
}

type Scan struct {
	events.Common
	SystemAddress         uint64
	StarSystem            string
	ScanType              string
	BodyID                int
	BodyName              string
	Parents               []ggja.BareObj
	DistanceFromArrivalLS float64
	Landable              bool
	Materials             []ScanMaterial
	ReserveLevel          string
	WasDiscovered         bool
	WasMapped             bool
}

func init() {
	events.RegisterType(string(ScanEvent), ScanEvent)
}
