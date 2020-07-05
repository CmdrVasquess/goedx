package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CmdrVasquess/goedx"
	bolt "go.etcd.io/bbolt"
)

func dumpfile(name string) {
	db, err := bolt.Open(name, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("systems"))
		crsr := b.Cursor()
		var sys goedx.System
		for k, v := crsr.First(); k != nil; k, v = crsr.Next() {
			dec := gob.NewDecoder(bytes.NewReader(v))
			if err = dec.Decode(&sys); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%d '%s' %v [%s / %s]\n",
				sys.Addr, sys.Name, sys.Coos,
				sys.FirstAccess.Format(time.RFC822),
				sys.LastAccess.Format(time.RFC822),
			)
		}
		return nil
	})
}

func main() {
	for _, arg := range os.Args[1:] {
		dumpfile(arg)
	}
}
