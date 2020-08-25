package bboltgalaxy

import (
	"strings"

	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/journal"
)

const UnknownType = 0

type System struct {
	goedx.System
	BodyCount int
	Bodies    []*Body
}

func (sys *System) GetBody(id int) (body *Body, newBody bool) {
	if id < len(sys.Bodies) {
		res := sys.Bodies[id]
		if res != nil {
			return res, false
		}
	} else {
		for id >= len(sys.Bodies) {
			sys.Bodies = append(sys.Bodies, nil)
		}
	}
	body = &Body{}
	sys.Bodies[id] = body
	return body, true
}

type BodyType int

const (
	Star BodyType = 1 + iota
	Belt
	BeltCluster
	Planet
	Barycenter
)

func BodyTypeFromScan(scan *journal.Scan) BodyType {
	switch {
	case scan.StarType != "":
		return Star
	case scan.PlanetClass != "":
		return Planet
	case strings.Index(scan.BodyName, "Belt Cluster") >= 0:
		return BeltCluster
	}
	return 0
}

type Ring struct {
	Name           string
	Class          string
	Mass           float32
	RadMin, RadMax float32
}

type Body struct {
	Parent int
	Name   string
	Type   BodyType
	DistFA float32            `json:",omitempty"`
	Mats   map[string]float32 `json:",omitempty"`
	Rings  []Ring             `json:",omitempty"`
}

func (b *Body) SetMat(mat string, f float32) {
	if b.Mats == nil {
		b.Mats = make(map[string]float32)
	}
	b.Mats[mat] = f
}
