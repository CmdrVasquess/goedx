package bboltgalaxy

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"strings"
	"time"

	"github.com/CmdrVasquess/goedx/journal"

	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/events"
	bolt "go.etcd.io/bbolt"
)

type Galaxy struct {
	db       *bolt.DB
	lastAddr uint64
	lastSys  *System
}

func Open(file string) (*Galaxy, error) {
	db, err := bolt.Open(file, 0666, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bktSystems)
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Galaxy{db: db}, nil
}

func (gxy *Galaxy) Close() error {
	return gxy.db.Close()
}

func (gxy *Galaxy) LastSystem() (uint64, *System) {
	return gxy.lastAddr, gxy.lastSys
}

func (gxy *Galaxy) EdgxSystem(
	addr uint64,
	name string,
	coos []float32,
	touch time.Time,
) (sys *goedx.System, tok interface{}) {
	res := gxy.FindSystemByAddr(addr)
	if res == nil {
		res = &System{
			System: *goedx.NewSystem(addr, name, coos...),
		}
		res.FirstAccess = touch
		res.LastAccess = touch
		gxy.UpdateSystem(res)
	} else if !res.Same(name, coos...) {
		res.Set(name, coos...)
		if !touch.IsZero() {
			res.LastAccess = touch
		}
		gxy.UpdateSystem(res)
	} else if !touch.IsZero() && res.LastAccess.Add(5*time.Minute).Before(touch) { // TODO param
		res.LastAccess = touch
		gxy.UpdateSystem(res)
	}
	return &res.System, res
}

func (gxy *Galaxy) FindSystemByAddr(addr uint64) (res *System) {
	if gxy.lastAddr == addr {
		return gxy.lastSys
	}
	err := gxy.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bktSystems)
		raw := b.Get(addr2key(addr))
		if raw == nil {
			return nil
		}
		dec := gob.NewDecoder(bytes.NewReader(raw))
		res = new(System)
		err := dec.Decode(res)
		if err != nil {
			res = nil
		}
		return err
	})
	if err == nil && addr != 0 {
		gxy.lastAddr = addr
		gxy.lastSys = res
	}
	return res
}

func (gxy *Galaxy) UpdateSystem(sys *System) {
	gxy.db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket(bktSystems)
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if err = enc.Encode(sys); err == nil {
			err = b.Put(addr2key(sys.Addr), buf.Bytes())
		}
		return err
	})
}

var bktSystems = []byte("systems")

func addr2key(addr uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, addr)
	return res
}

func (gxy *Galaxy) PrepareEDEvent(e events.Event) (token interface{}) {
	switch evt := e.(type) {
	case *journal.Scan:
		if evt.SystemAddress == 0 {
			return nil
		}
		return journal.ScanEvent
	case *journal.FSSDiscoveryScan:
		return journal.FSSDiscoveryScanEvent
	}
	return nil
}

func (gxy *Galaxy) FinishEDEvent(token interface{}, e events.Event, _ goedx.Change) {
	switch evt := e.(type) {
	case *journal.Scan:
		gxy.evtScan(evt)
	case *journal.FSSDiscoveryScan:
		gxy.evtFSSDisco(evt)
	}
}

func (gxy *Galaxy) evtFSSDisco(scan *journal.FSSDiscoveryScan) {
	if scan.BodyCount == 0 {
		return
	}
	chg := false
	sys := gxy.FindSystemByAddr(scan.SystemAddress)
	if sys == nil {
		sys = &System{
			System: *goedx.NewSystem(
				scan.SystemAddress,
				scan.SystemName,
			),
			BodyCount: scan.BodyCount,
		}
		chg = true
	} else if scan.BodyCount != sys.BodyCount {
		sys.BodyCount = scan.BodyCount
		chg = true
	}
	if chg {
		gxy.UpdateSystem(sys)
	}
}

func (gxy *Galaxy) evtScan(scan *journal.Scan) {
	sys := gxy.FindSystemByAddr(scan.SystemAddress)
	chg := false
	if sys == nil {
		sys = &System{
			System: *goedx.NewSystem(
				scan.SystemAddress,
				scan.StarSystem,
			),
		}
		chg = true
	}
	body, newBody := sys.GetBody(scan.BodyID)
	chg = chg || newBody
	if body.Name != scan.BodyName {
		body.Name = scan.BodyName
		chg = true
	}
	if scan.DistanceFromArrivalLS > 0 {
		body.DistFA = float32(scan.DistanceFromArrivalLS)
		chg = true
	}
	if t := BodyTypeFromScan(scan); t != body.Type {
		body.Type = t
		chg = true
	}
	parents := parents(scan)
	chg = linkParents(sys, body, parents) || chg
	if len(parents) == 0 {
		body.Parent = -1
	} else {
		body.Parent = parents[0].id
	}
	if body.Type == BeltCluster {
		chg = scanBeltCluster(scan, sys, body, parents) || chg
	}
	chg = scanRings(scan, body) || chg
	chg = scanMats(scan, body) || chg
	if chg {
		gxy.UpdateSystem(sys)
	}
}

func linkParents(sys *System, body *Body, ps []parent) (chg bool) {
	for _, p := range ps {
		if body.Parent != p.id {
			body.Parent = p.id
			chg = true
		}
		pb, newPb := sys.GetBody(p.id)
		chg = chg || newPb
		if pb.Type == 0 && p.typ != 0 {
			pb.Type = p.typ
			chg = true
		}
		body = pb
	}
	return chg
}

func ringClass(scanClass string) string {
	const prefix = "eRingClass_"
	if strings.HasPrefix(scanClass, prefix) {
		return scanClass[len(prefix):]
	}
	return scanClass
}

func scanRings(scan *journal.Scan, body *Body) (chg bool) {
	if len(scan.Rings) == 0 {
		return false
	}
	if len(body.Rings) < len(scan.Rings) {
		body.Rings = make([]Ring, len(scan.Rings))
		for i, sr := range scan.Rings {
			body.Rings[i] = Ring{
				Name:   sr.Name,
				Class:  ringClass(sr.RingClass),
				Mass:   float32(sr.MassMT),
				RadMin: float32(sr.InnerRad),
				RadMax: float32(sr.OuterRad),
			}
		}
		chg = true
	}
	return chg
}

func scanMats(scan *journal.Scan, body *Body) (chg bool) {
	if len(scan.Materials) == 0 {
		return false
	}
	for _, mat := range scan.Materials {
		body.SetMat(mat.Name, mat.Percent/100)
	}
	return true
}

func beltName(clusterName string) string {
	pos := strings.Index(clusterName, "Belt Cluster")
	if pos < 0 {
		return "Belt"
	}
	return clusterName[:pos+4]
}

func scanBeltCluster(scan *journal.Scan, sys *System, body *Body, ps []parent) (chg bool) {
	if len(ps) == 0 {
		return false
	}
	beltId := ps[0].id
	belt, chg := sys.GetBody(beltId)
	if belt.Type != Belt {
		belt.Type = Belt
		chg = true
	}
	if belt.Name == "" {
		belt.Name = beltName(body.Name)
		chg = true
	}
	return chg
}

type parent struct {
	typ BodyType
	id  int
}

var parentTypes = map[string]BodyType{
	"Null": Barycenter,
	"Star": Star,
	"Ring": Belt,
}

func parents(scan *journal.Scan) []parent {
	if len(scan.Parents) == 0 {
		return nil
	}
	res := make([]parent, len(scan.Parents))
	for i, p := range scan.Parents {
		r := &res[i]
		for k, v := range p { // len(p) == 1
			r.typ = parentTypes[k]
			r.id = int(v.(float64))
		}
	}
	return res
}
