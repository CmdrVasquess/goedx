package goedx

type Galaxy interface {
	EdgxSystem(addr uint64, name string, coos []float32) (*System, interface{})
}

const EchoGalaxy = echoGxy(0)

type echoGxy int

func (_ echoGxy) EdgxSystem(addr uint64, name string, coos []float32) (*System, interface{}) {
	res := NewSystem(addr, name, coos...)
	return res, nil
}

type InMemGalaxy map[uint64]*System

func (g InMemGalaxy) EdgxSystem(addr uint64, name string, coos []float32) (*System, interface{}) {
	res := g[addr]
	if res == nil {
		res = NewSystem(addr, name, coos...)
	} else {
		l := len(coos)
		if l > 3 {
			l = 3
		}
		for l--; l >= 0; l-- {
			res.Coos[l].Set(coos[l], 0)
		}
	}
	return res, nil
}
