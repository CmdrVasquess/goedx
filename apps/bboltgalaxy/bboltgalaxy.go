package bboltgalaxy

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"time"

	"github.com/CmdrVasquess/goedx"
	bolt "go.etcd.io/bbolt"
)

type Galaxy bolt.DB

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
	return (*Galaxy)(db), nil
}

func (g *Galaxy) Close() error {
	return (*bolt.DB)(g).Close()
}

func (g *Galaxy) EdgxSystem(
	addr uint64,
	name string,
	coos []float32,
	touch time.Time,
) (sys *goedx.System, tok interface{}) {
	var res *goedx.System
	db := (*bolt.DB)(g)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bktSystems)
		raw := b.Get(addr2key(addr))
		if raw == nil {
			return nil
		}
		dec := gob.NewDecoder(bytes.NewReader(raw))
		res = new(goedx.System)
		err := dec.Decode(res)
		if err != nil {
			res = nil
		}
		return err
	})
	if res == nil {
		db.Update(func(tx *bolt.Tx) (err error) {
			b := tx.Bucket(bktSystems)
			res = goedx.NewSystem(addr, name, coos...)
			res.FirstAccess = touch
			res.LastAccess = touch
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err = enc.Encode(res); err == nil {
				err = b.Put(addr2key(addr), buf.Bytes())
			}
			return err
		})
	} else if !res.Same(name, coos...) {
		res.Set(name, coos...)
		if !touch.IsZero() {
			res.LastAccess = touch
		}
		g.UpdateSystem(res)
	} else if !touch.IsZero() && res.LastAccess.Add(5*time.Minute).Before(touch) { // TODO param
		res.LastAccess = touch
		g.UpdateSystem(res)
	}
	return res, nil
}

func (g *Galaxy) UpdateSystem(sys *goedx.System) {
	db := (*bolt.DB)(g)
	db.Update(func(tx *bolt.Tx) (err error) {
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
