package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	fDB     string
	evtCmdr = []byte(`"event":"Commander"`)
	evtLdg  = []byte(`"event":"LoadGame"`)
	evtFSDJ = []byte(`"event":"FSDJump"`)
)

type LoadGame struct {
	Commander string
}

type Commander struct {
	FID  string
	Name string
}

type FSDJump struct {
	Timestamp     time.Time
	StarSystem    string
	SystemAddress int64
	StarPos       [3]float64
}

var cmdrId int

func importFrom(db *sql.DB, rd io.Reader) {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if p := recover(); p == nil {
			tx.Commit()
		} else {
			log.Println("rollback transaction")
			tx.Rollback()
		}
	}()
	scn := bufio.NewScanner(rd)
	for scn.Scan() {
		line := scn.Bytes()
		switch {
		case bytes.Index(line, evtFSDJ) >= 0:
			if cmdrId <= 0 {
				log.Panicf("no commander for jump: %s", string(line))
			}
			fsdJump(tx, line, cmdrId)
		case bytes.Index(line, evtCmdr) >= 0:
			cmdrId = cmdrEvent(tx, line)
			if cmdrId <= 0 {
				log.Panicf("no commander for jump: %s", string(line))
			}
		case bytes.Index(line, evtLdg) >= 0:
			cmdrId = ldgEvent(tx, line)
			if cmdrId <= 0 {
				log.Panicf("no commander for jump: %s", string(line))
			}
		}
	}
}

func switchCmdr(db *sql.Tx, cmdrNm, cmdrFid string) int {
	if cmdrNm == "" {
		log.Panic("commander without name")
	}
	var cmdrId int
	var fid sql.NullString
	err := db.QueryRow(`SELECT id, fid FROM cmdrs WHERE name=$1`, cmdrNm).
		Scan(&cmdrId, &fid)
	switch {
	case err == sql.ErrNoRows:
		var res sql.Result
		var err error
		if cmdrFid == "" {
			res, err = db.Exec(`INSERT INTO cmdrs (name) VALUES ($1)`, cmdrNm)
		} else {
			res, err = db.Exec(`INSERT INTO cmdrs (fid, name) VALUES ($1, $2)`,
				cmdrFid, cmdrNm)
		}
		if err != nil {
			log.Panic(err)
		}
		if id, err := res.LastInsertId(); err != nil {
			log.Panic(err)
		} else {
			cmdrId = int(id)
		}
		return cmdrId
	case err != nil:
		log.Panic(err)
	}
	if (!fid.Valid || fid.String == "") && cmdrFid != "" {
		_, err = db.Exec(`UPDATE cmdrs SET fid=$1 WHERE id=$2`, cmdrFid, cmdrId)
		if err != nil {
			log.Printf("cannot update fid %s of commander %s: %s",
				cmdrFid, cmdrNm, err)
		}
	}
	return cmdrId
}

func cmdrEvent(db *sql.Tx, line []byte) int {
	var cmdr Commander
	if err := json.Unmarshal(line, &cmdr); err != nil {
		log.Panic(err)
	}
	return switchCmdr(db, cmdr.Name, cmdr.FID)
}

func ldgEvent(db *sql.Tx, line []byte) int {
	var ldg LoadGame
	if err := json.Unmarshal(line, &ldg); err != nil {
		log.Panic(err)
	}
	return switchCmdr(db, ldg.Commander, "")
}

func fsdJump(db *sql.Tx, line []byte, cmdrId int) {
	var jump FSDJump
	if err := json.Unmarshal(line, &jump); err != nil {
		log.Panic(err)
	}
	sysId := getSystem(db, &jump)
	if sysId <= 0 {
		log.Panicf("no system for jump: %+v", &jump)
	}
	_, err := db.Exec(`INSERT INTO visits (cmdr, sys, arrive) VALUES ($1, $2, $3)`,
		cmdrId, sysId, jump.Timestamp)
	if err != nil {
		log.Panic(err)
	}
}

func getSystem(db *sql.Tx, jump *FSDJump) int {
	var sysId int
	if jump.SystemAddress != 0 {
		err := db.QueryRow(`SELECT id FROM systems WHERE addr = $1`,
			jump.SystemAddress).
			Scan(&sysId)
		switch {
		case err == nil:
			return sysId
		case err != sql.ErrNoRows:
			log.Panic(err)
		}
	}
	var addr sql.NullInt64
	err := db.QueryRow(`SELECT id, addr FROM systems WHERE lower(name)=lower($1)`,
		jump.StarSystem).
		Scan(&sysId, &addr)
	switch {
	case err == sql.ErrNoRows:
		if jump.SystemAddress > 0 {
			res, err := db.Exec(
				`INSERT INTO systems (name, addr, x, y, z) VALUES ($1, $2, $3, $4, $5)`,
				jump.StarSystem,
				jump.SystemAddress,
				jump.StarPos[0], jump.StarPos[1], jump.StarPos[2])
			if err != nil {
				log.Panic(err)
			}
			if id, err := res.LastInsertId(); err != nil {
				log.Panic(err)
			} else {
				return int(id)
			}
		} else {
			res, err := db.Exec(
				`INSERT INTO systems (name, x, y, z) VALUES ($1, $2, $3, $4)`,
				jump.StarSystem,
				jump.StarPos[0], jump.StarPos[1], jump.StarPos[2])
			if err != nil {
				log.Panic(err)
			}
			if id, err := res.LastInsertId(); err != nil {
				log.Panic(err)
			} else {
				return int(id)
			}
		}
	case err != nil:
		log.Panic(err)
	}
	if !addr.Valid || addr.Int64 == 0 {
		if jump.SystemAddress != 0 {
			_, err := db.Exec(`UPDATE systems SET addr=$1 WHERE id=$2`,
				jump.SystemAddress,
				sysId)
			if err != nil {
				log.Printf("failed to set system addr %d for %s: %s",
					jump.SystemAddress, jump.StarSystem, err)
			}
		}
	} else if jump.SystemAddress != 0 && addr.Int64 != jump.SystemAddress {
		log.Println("ambg addr for %s: %d / %d",
			jump.StarSystem, addr, jump.SystemAddress)
	}
	return sysId
}

func importLog(db *sql.DB, file string) {
	rd, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer rd.Close()
	importFrom(db, rd)
}

func importGz(db *sql.DB, file string) {
	rd, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer rd.Close()
	rdgz, err := gzip.NewReader(rd)
	if err != nil {
		log.Panic(err)
	}
	defer rdgz.Close()
	importFrom(db, rdgz)
}

func main() {
	flag.StringVar(&fDB, "db", "", "sqlite3 DB file")
	flag.Parse()
	db, err := sql.Open("sqlite3", fDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for _, arg := range flag.Args() {
		log.Printf("import %s", arg)
		switch filepath.Ext(arg) {
		case ".log":
			importLog(db, arg)
		case ".gz":
			importGz(db, arg)
		}
	}
}
