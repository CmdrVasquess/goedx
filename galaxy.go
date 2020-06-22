package edgx

type Galaxy interface {
	EdgxSystem(addr uint64, name string, coos []float32) (*System, interface{})
}

const EchoGalaxy = echoGxy(0)

type echoGxy int

func (_ echoGxy) EdgxSystem(addr uint64, name string, coos []float32) (*System, interface{}) {
	res := &System{
		Addr: addr,
		Name: name,
	}
	l := len(coos)
	if l > 3 {
		l = 3
	}
	for l--; l >= 0; l-- {
		res.Coos[l].Set(coos[l], 0)
	}
	return res, nil
}
