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

type System struct {
	goedx.System
	FirstAccess time.Time
	LastAccess  time.Time
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
	return (*Galaxy)(db), nil
}

func (g *Galaxy) Close() error {
	return (*bolt.DB)(g).Close()
}

func (g *Galaxy) EdgxSystem(
	addr uint64,
	name string,
	coos []float32,
) (sys *goedx.System, tok interface{}) {
	var res *System
	db := (*bolt.DB)(g)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bktSystems)
		raw := b.Get(addr2key(addr))
		if raw == nil {
			return nil
		}
		dec := gob.NewDecoder(bytes.NewReader(raw))
		res := new(System)
		err := dec.Decode(res)
		if err != nil {
			res = nil
		}
		return err
	})
	now := time.Now()
	if res == nil {
		db.Update(func(tx *bolt.Tx) (err error) {
			b := tx.Bucket(bktSystems)
			res := System{
				System:      *goedx.NewSystem(addr, name, coos...),
				FirstAccess: now,
				LastAccess:  now,
			}
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err = enc.Encode(res); err == nil {
				err = b.Put(addr2key(addr), buf.Bytes())
			}
			return err
		})
	} else if !res.Same(name, coos...) {
		res.Set(name, coos...)
		res.LastAccess = now
		g.UpdateSystem(res)
	} else if res.LastAccess.Add(5 * time.Minute).Before(now) { // TODO param
		res.LastAccess = now
		g.UpdateSystem(res)
	}
	return &res.System, nil
}

func (g *Galaxy) UpdateSystem(sys *System) {
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
