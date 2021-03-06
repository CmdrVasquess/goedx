package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"

	"github.com/CmdrVasquess/goedx/apps/bboltgalaxy"

	bolt "go.etcd.io/bbolt"
)

func dumpfile(name string) {
	db, err := bolt.Open(name, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	jenc := json.NewEncoder(os.Stdout)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("systems"))
		crsr := b.Cursor()
		var sys bboltgalaxy.System
		for k, v := crsr.First(); k != nil; k, v = crsr.Next() {
			dec := gob.NewDecoder(bytes.NewReader(v))
			if err = dec.Decode(&sys); err != nil {
				log.Fatal(err)
			}
			jenc.Encode(&sys)
		}
		return nil
	})
}

func main() {
	for _, arg := range os.Args[1:] {
		dumpfile(arg)
	}
}
