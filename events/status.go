package events

type StatusFlag uint32

type statusT string

const StatusEvent = statusT("Status")

func (t statusT) New() Event     { return new(Status) }
func (t statusT) String() string { return string(t) }

type Status struct {
	Common
	Flags     StatusFlag
	Pips      [3]int
	FireGroup int
	Fuel      struct {
		Main      float64 `json:"FuelMain"`
		Reservoir float64 `json:"FuelReservoi"`
	}
}

func (s *Status) AnyFlag(fs StatusFlag) bool {
	return s.Flags&fs > 0
}

func (s *Status) AllFlags(fs StatusFlag) bool {
	return s.Flags&fs == fs
}

func init() {
	RegisterType(string(StatusEvent), StatusEvent)
}
