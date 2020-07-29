package goedx

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"git.fractalqb.de/fractalqb/ggja"
)

const (
	LocTypeSystem  = "system"
	LocTypePort    = "port"
	LocTypeFSDJump = "jump"

	jsonTypeTag = "@type"
	jsonSysTag  = "Sys"
)

type Location interface {
	System() *System
	ToMap(m *map[string]interface{}, setType bool) error
	FromMap(m map[string]interface{}) error
}

type System struct {
	Addr        uint64
	Name        string
	Coos        SysCoos
	FirstAccess time.Time
	LastAccess  time.Time
}

func NewSystem(addr uint64, name string, coos ...float32) *System {
	res := &System{Addr: addr}
	res.Set(name, coos...)
	return res
}

func (s *System) Set(name string, coos ...float32) {
	if name != "" {
		s.Name = name
	}
	l := len(coos)
	if l > 3 {
		l = 3
	}
	for l--; l >= 0; l-- {
		s.Coos[l] = ChgF32(coos[l])
	}
}

func (s *System) Same(name string, coos ...float32) bool {
	if name != "" && s.Name != name {
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

func (s *System) ToMap(m *map[string]interface{}, setType bool) error {
	(*m)["Addr"] = s.Addr
	(*m)["Name"] = s.Name
	if s.Coos.Valid() {
		(*m)["Coos"] = &s.Coos
	}
	if !s.FirstAccess.IsZero() {
		(*m)["FirstAccess"] = s.FirstAccess
	}
	if !s.LastAccess.IsZero() {
		(*m)["LastAccess"] = s.LastAccess
	}
	if setType {
		(*m)[jsonTypeTag] = LocTypeSystem
	}
	return nil
}

func (s *System) FromMap(m map[string]interface{}) (err error) {
	obj := ggja.Obj{Bare: m, OnError: func(e error) { err = e }}
	s.Addr = obj.MUint64("Addr")
	if err != nil {
		return err
	}
	s.Name = obj.MStr("Name")
	if err != nil {
		return err
	}
	if coos := obj.Arr("Coos"); coos == nil {
		nan32 := ChgF32(math.NaN())
		s.Coos[0] = nan32
		s.Coos[1] = nan32
		s.Coos[2] = nan32
	}
	if ts := obj.Str("FirstAccess", ""); ts == "" {
		s.FirstAccess = time.Time{}
	} else {
		data := []byte{'"'}
		data = append(data, []byte(ts)...)
		data = append(data, '"')
		if err := json.Unmarshal(data, &s.FirstAccess); err != nil {
			return err
		}
	}
	if ts := obj.Str("LastAccess", ""); ts == "" {
		s.LastAccess = time.Time{}
	} else {
		data := []byte{'"'}
		data = append(data, []byte(ts)...)
		data = append(data, '"')
		if err := json.Unmarshal(data, &s.LastAccess); err != nil {
			return err
		}
	}
	return nil
}

type Port struct {
	Sys    *System
	Name   string
	Type   string `json:",omitempty"`
	Docked bool
}

func (p *Port) System() *System { return p.Sys }

func (p *Port) ToMap(m *map[string]interface{}, setType bool) error {
	sys := make(map[string]interface{})
	if err := p.Sys.ToMap(&sys, false); err != nil {
		return fmt.Errorf("Port.Sys: %s", err)
	}
	(*m)[jsonSysTag] = sys
	(*m)["Name"] = p.Name
	if p.Type != "" {
		(*m)["Type"] = p.Type
	}
	(*m)["Docked"] = p.Docked
	if setType {
		(*m)[jsonTypeTag] = LocTypePort
	}
	return nil
}

func (p *Port) FromMap(m map[string]interface{}) (err error) {
	obj := ggja.Obj{Bare: m, OnError: func(e error) { err = e }}
	sysMap := obj.MObj(jsonSysTag)
	if err != nil {
		return err
	}
	sys := new(System)
	if err = sys.FromMap(sysMap.Bare); err != nil {
		return fmt.Errorf("Port.Sys: %s", err)
	}
	p.Sys = sys
	p.Name = obj.MStr("Name")
	if err != nil {
		return err
	}
	p.Type = obj.Str("Type", "")
	if err != nil {
		return err
	}
	p.Docked = obj.MBool("Docked")
	return err
}

type FSDJump struct {
	Sys *System
	To  *System
}

func (j *FSDJump) System() *System { return j.Sys }

func (j *FSDJump) ToMap(m *map[string]interface{}, setType bool) error {
	sys := make(map[string]interface{})
	if err := j.Sys.ToMap(&sys, false); err != nil {
		return err
	}
	to := make(map[string]interface{})
	if err := j.To.ToMap(&to, false); err != nil {
		return err
	}
	(*m)[jsonSysTag] = sys
	(*m)["To"] = to
	if setType {
		(*m)[jsonTypeTag] = LocTypeFSDJump
	}
	return nil
}

func (j *FSDJump) FromMap(m map[string]interface{}) (err error) {
	obj := ggja.Obj{Bare: m, OnError: func(e error) { err = e }}
	sysMap := obj.MObj(jsonSysTag)
	if err != nil {
		return err
	}
	sys := new(System)
	if err = sys.FromMap(sysMap.Bare); err != nil {
		return err
	}
	j.Sys = sys
	sysMap = obj.MObj("To")
	if err != nil {
		return err
	}
	sys = new(System)
	if err = sys.FromMap(sysMap.Bare); err != nil {
		return err
	}
	j.To = sys
	return nil
}

type JSONLocation struct {
	Location
}

func (jl JSONLocation) System() *System {
	if jl.Location == nil {
		return nil
	}
	return jl.Location.System()
}

func (jl JSONLocation) FSDJump() *FSDJump {
	if j, ok := jl.Location.(*FSDJump); ok {
		return j
	}
	return nil
}

func (jl JSONLocation) Port() *Port {
	if p, ok := jl.Location.(*Port); ok {
		return p
	}
	return nil
}

var jsonNull = []byte("null")

func (jloc JSONLocation) MarshalJSON() ([]byte, error) {
	if jloc.Location == nil {
		return jsonNull, nil
	}
	tmp := make(map[string]interface{})
	err := jloc.Location.ToMap(&tmp, true)
	if err != nil {
		return nil, err
	}
	return json.Marshal(tmp)
}

func (jloc *JSONLocation) UnmarshalJSON(data []byte) (err error) {
	tmp := make(ggja.BareObj)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	if len(tmp) == 0 {
		jloc.Location = nil
		return nil
	}
	obj := ggja.Obj{Bare: tmp, OnError: func(e error) { err = e }}
	switch obj.Str(jsonTypeTag, "") {
	case LocTypeSystem:
		s := new(System)
		if err := s.FromMap(tmp); err != nil {
			return err
		}
		jloc.Location = s
	case LocTypePort:
		p := new(Port)
		if err := p.FromMap(tmp); err != nil {
			return err
		}
		jloc.Location = p
	case LocTypeFSDJump:
		j := new(FSDJump)
		if err := j.FromMap(tmp); err != nil {
			return err
		}
		jloc.Location = j
	case "":
		err = errors.New("missing @type attribute")
	default:
		err = fmt.Errorf("unkown location type '%s'", tmp["@type"])
	}
	return err
}
