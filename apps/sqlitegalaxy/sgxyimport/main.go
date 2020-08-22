package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const sqliteTs = "2006-01-02 15:04:05"

var (
	fDB          string
	evtCmdr      = []byte(`"event":"Commander"`)
	evtLdg       = []byte(`"event":"LoadGame"`)
	evtFSDJ      = []byte(`"event":"FSDJump"`)
	evtDocked    = []byte(`"event":"Docked"`)
	startAfter   time.Time
	countSystems int
	countPorts   int
	countVisits  int
	countDocked  int
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

type Docked struct {
	Timestamp   time.Time
	StarSystem  string
	StationName string
	StationType string
}

var (
	cmdrId        int
	currentSystem int
)

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
		case bytes.Index(line, evtDocked) >= 0:
			docked(tx, line)
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
	currentSystem = 0
	return switchCmdr(db, ldg.Commander, "")
}

func fsdJump(db *sql.Tx, line []byte, cmdrId int) {
	var jump FSDJump
	if err := json.Unmarshal(line, &jump); err != nil {
		log.Panic(err)
	}
	if !jump.Timestamp.After(startAfter) {
		return
	}
	currentSystem = getSystem(db, jump.StarSystem, jump.SystemAddress, &jump.StarPos)
	if currentSystem <= 0 {
		log.Panicf("no system for jump: %+v", &jump)
	}
	_, err := db.Exec(`INSERT INTO visits (cmdr, sys, arrive) VALUES ($1, $2, $3)`,
		cmdrId, currentSystem, jump.Timestamp)
	if err != nil {
		log.Panic(err)
	}
	countVisits++
}

func docked(db *sql.Tx, line []byte) {
	var dock Docked
	if err := json.Unmarshal(line, &dock); err != nil {
		log.Panic(err)
	}
	if dock.StationName == "" {
		log.Println("no port name in ", string(line))
		return
	}
	currentSystem = getSystem(db, dock.StarSystem, 0, nil)
	pid := getPort(db, currentSystem, dock.StationName, dock.StationType)
	_, err := db.Exec(`INSERT INTO docked (cmdr, port, arrive) VALUES ($1, $2, $3)`,
		cmdrId, pid, dock.Timestamp)
	if err != nil {
		log.Panic(err)
	}
	countDocked++
}

func getPort(db *sql.Tx, sys int, name, typ string) (pid int) {
	if sys == 0 {
		log.Printf("searching port %s (%s) in system 0", name, typ)
		return 0
	}
	err := db.QueryRow(`SELECT id FROM ports WHERE sys=$1 and name=$2`,
		sys,
		name,
	).Scan(&pid)
	if err == nil {
		return pid
	}
	res, err := db.Exec(`INSERT INTO ports (sys, name, type) VALUES ($1, $2, $3)`,
		sys, name, strings.ToLower(typ))
	if err != nil {
		log.Panic(err)
	}
	if id, err := res.LastInsertId(); err != nil {
		log.Panic(err)
	} else {
		pid = int(id)
	}
	countPorts++
	return pid
}

func sysAddCoos(db *sql.Tx, sysId int, starPos *[3]float64) {
	_, err := db.Exec(`UPDATE systems SET x=$1, y=$2, z=$3 WHERE id=$4`,
		starPos[0], starPos[1], starPos[2],
		sysId)
	if err != nil {
		log.Printf("failed to set system coos for %d", sysId)
	}
}

func getSystem(db *sql.Tx, sysName string, sysAddr int64, starPos *[3]float64) int {
	var sysId int
	var sysx, sysy, sysz sql.NullFloat64
	if sysAddr != 0 {
		err := db.QueryRow(`SELECT id, x FROM systems WHERE addr = $1`, sysAddr).
			Scan(&sysId, &sysx)
		switch {
		case err == nil:
			if !sysx.Valid && starPos != nil {
				sysAddCoos(db, sysId, starPos)
			}
			return sysId
		case err != sql.ErrNoRows:
			log.Panic(err)
		}
	}
	var addr sql.NullInt64
	err := db.QueryRow(`SELECT id, addr FROM systems WHERE lower(name)=lower($1)`,
		sysName).
		Scan(&sysId, &addr)
	switch {
	case err == sql.ErrNoRows:
		if starPos != nil {
			sysx.Valid = true
			sysx.Float64 = starPos[0]
			sysy.Valid = true
			sysy.Float64 = starPos[1]
			sysz.Valid = true
			sysz.Float64 = starPos[2]
		}
		if sysAddr > 0 {
			res, err := db.Exec(
				`INSERT INTO systems (name, addr, x, y, z) VALUES ($1, $2, $3, $4, $5)`,
				sysName,
				sysAddr,
				sysx, sysy, sysz)
			if err != nil {
				log.Panic(err)
			}
			countSystems++
			if id, err := res.LastInsertId(); err != nil {
				log.Panic(err)
			} else {
				return int(id)
			}
		} else {
			res, err := db.Exec(
				`INSERT INTO systems (name, x, y, z) VALUES ($1, $2, $3, $4)`,
				sysName,
				sysx, sysy, sysz)
			if err != nil {
				log.Panic(err)
			}
			countSystems++
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
		if sysAddr != 0 {
			_, err := db.Exec(`UPDATE systems SET addr=$1 WHERE id=$2`,
				sysAddr,
				sysId)
			if err != nil {
				log.Printf("failed to set system addr %d for %s: %s",
					sysAddr, sysName, err)
			}
		}
	} else if sysAddr != 0 && addr.Int64 != sysAddr {
		log.Println("ambg addr for %s: %d / %d",
			sysName, addr, sysAddr)
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

func lastJump(db *sql.DB) (t time.Time, err error) {
	var tstr sql.NullString
	if err = db.QueryRow(`SELECT max(arrive) FROM visits`).Scan(&tstr); err != nil {
		return t, err
	}
	if !tstr.Valid {
		return time.Time{}, nil
	}
	t, err = time.Parse(sqliteTs, tstr.String[:len(sqliteTs)])
	return t, err
}

func main() {
	flag.StringVar(&fDB, "db", "", "sqlite3 DB file")
	flag.Parse()
	sigs := make(chan os.Signal, 1) // '1' is important for select to not always default
	signal.Notify(sigs, os.Interrupt)
	db, err := sql.Open("sqlite3", fDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	startAfter, err = lastJump(db)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
ARGS_LOOP:
	for _, arg := range flag.Args() {
		log.Printf("import %s", arg)
		switch filepath.Ext(arg) {
		case ".log":
			importLog(db, arg)
		case ".gz":
			importGz(db, arg)
		}
		select {
		case <-sigs:
			log.Println("Import interrupted")
			break ARGS_LOOP
		default:
		}
	}
	fmt.Printf(`Imported since %s:
- %d new system
- %d new visits
- %d new ports
- %d new dockings
`,
		startAfter,
		countSystems, countVisits,
		countPorts, countDocked)
}
