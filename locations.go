package edgx

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.fractalqb.de/fractalqb/ggja"
)

type Location interface {
	System() *System
}

type System struct {
	Addr uint64
	Name string
	Coos []float32
}

func (s *System) System() *System { return s }

type Port struct {
	Sys  *System
	Name string
}

func (p *Port) System() *System { return p.Sys }

type JSONLocation struct {
	Location
}

func (jloc JSONLocation) MarshalJSON() ([]byte, error) {
	t := map[string]interface{}{
		"v": jloc.Location,
	}
	switch jloc.Location.(type) {
	case *Port:
		t["t"] = "port"
	default:
		return nil, fmt.Errorf("unknown location type '%T'", jloc.Location)
	}
	return json.Marshal(&t)
}

func fillSys(sys *System, obj *ggja.Obj) {
	sys.Addr = obj.MUint64("Addr")
	sys.Name = obj.MStr("Name")
	for _, coo := range obj.MArr("Coos").Bare {
		sys.Coos = append(sys.Coos, float32(coo.(float64)))
	}
}

func (jloc *JSONLocation) UnmarshalJSON(data []byte) (err error) {
	tmp := make(ggja.GenObj)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	obj := ggja.Obj{Bare: tmp, OnError: func(e error) { err = e }}
	loc := obj.MObj("v")
	if err != nil {
		return err
	}
	switch obj.Str("t", "") {
	case "port":
		p := &Port{Name: loc.MStr("Name"), Sys: new(System)}
		if err != nil {
			return err
		}
		fillSys(p.Sys, loc.MObj("Sys"))
		jloc.Location = p
	case "":
		err = errors.New("missing @type attribute")
	default:
		err = fmt.Errorf("unkown location type '%s'", tmp["@type"])
	}
	return err
}
