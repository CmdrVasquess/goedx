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

type ScanRing struct {
	Name      string
	RingClass string
	MassMT    float64
	InnerRad  float64
	OuterRad  float64
}

type Scan struct {
	events.Common
	SystemAddress         uint64
	StarSystem            string
	ScanType              string
	StarType              string
	PlanetClass           string
	BodyID                int
	BodyName              string
	Parents               []ggja.BareObj
	DistanceFromArrivalLS float64
	Landable              bool
	Materials             []ScanMaterial
	ReserveLevel          string
	Rings                 []ScanRing
	WasDiscovered         bool
	WasMapped             bool
}

func (_ *Scan) EventType() events.Type { return ScanEvent }

func init() {
	events.RegisterType(string(ScanEvent), ScanEvent)
}
