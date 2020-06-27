package goedx

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.fractalqb.de/fractalqb/ggja"
	"github.com/mitchellh/mapstructure"
)

type Location interface {
	System() *System
}

type System struct {
	Addr uint64
	Name string
	Coos SysCoos
}

func NewSystem(addr uint64, name string, coos ...float32) *System {
	res := &System{Addr: addr}
	res.Set(name, coos...)
	return res
}

func (s *System) Set(name string, coos ...float32) {
	s.Name = name
	l := len(coos)
	if l > 3 {
		l = 3
	}
	for l--; l >= 0; l-- {
		s.Coos[l] = ChgF32(coos[l])
	}
}

func (s *System) Same(name string, coos ...float32) bool {
	if s.Name != name {
		return false
	}
	if len(coos) != 3 {
		return false
	}
	for i, r := range coos {
		if s.Coos[i] != ChgF32(r) {
			return false
		}
	}
	return true
}

func (s *System) System() *System { return s }

type Port struct {
	Sys    *System
	Name   string
	Docked bool
}

func (p *Port) System() *System { return p.Sys }

type JSONLocation struct {
	Location
}

func (jl JSONLocation) System() *System {
	if jl.Location == nil {
		return nil
	}
	return jl.Location.System()
}

func (jl JSONLocation) Port() *Port {
	if p, ok := jl.Location.(*Port); ok {
		return p
	}
	return nil
}

const jsonTypeTag = "@type"

var jsonNull = []byte("null")

func (jloc JSONLocation) MarshalJSON() ([]byte, error) {
	if jloc.Location == nil {
		return jsonNull, nil
	}
	tmp := make(map[string]interface{})
	err := mapstructure.Decode(jloc.Location, &tmp)
	if err != nil {
		return nil, err
	}
	switch jloc.Location.(type) {
	case *System:
		tmp[jsonTypeTag] = "system"
	case *Port:
		tmp[jsonTypeTag] = "port"
	default:
		return nil, fmt.Errorf("unknown location type '%T'", jloc.Location)
	}
	return json.Marshal(tmp)
}

func (jloc *JSONLocation) UnmarshalJSON(data []byte) (err error) {
	tmp := make(ggja.GenObj)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	obj := ggja.Obj{Bare: tmp, OnError: func(e error) { err = e }}
	switch obj.Str(jsonTypeTag, "") {
	case "system":
		s := new(System)
		if err := mapstructure.Decode(tmp, s); err != nil {
			return err
		}
		jloc.Location = s
	case "port":
		p := new(Port)
		if err := mapstructure.Decode(tmp, p); err != nil {
			return err
		}
		jloc.Location = p
	case "":
		err = errors.New("missing @type attribute")
	default:
		err = fmt.Errorf("unkown location type '%s'", tmp["@type"])
	}
	return err
}
