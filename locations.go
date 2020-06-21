package edgx

type Location interface {
	SystemAddr() uint64
	SystemName() string
}

type System struct {
	Addr uint64
	Name string
	Coos [3]float32
}

type Port struct {
	System *System
	Name   *string
}
