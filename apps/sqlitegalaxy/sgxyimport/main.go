package main

import (
	"bufio"
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

	"git.fractalqb.de/fractalqb/sqlize"
	"git.fractalqb.de/fractalqb/sqlize/null"
	sgx "github.com/CmdrVasquess/goedx/apps/sqlitegalaxy"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
	_ "github.com/mattn/go-sqlite3"
)

const sqliteTs = "2006-01-02 15:04:05"

type system struct {
	id   int
	addr int64
	name string
}

func (s *system) clear() {
	s.id = sqlize.NoID
	s.addr = sqlize.NoID
	s.name = ""
}

func (s *system) load(db sqlize.Tx, sysName string, sysAddr int64, starPos *[3]float64) {
	var sysx, sysy, sysz sql.NullFloat64
	if sysAddr != 0 {
		s.addr = sysAddr
		err := db.QueryRow(sgx.SqlSysByAddr.MustSQL(), sysAddr).
			Scan(&s.id, &s.name, &sysx, &sysy, &sysz)
		switch {
		case err == nil:
			if !sysx.Valid && starPos != nil {
				sysAddCoos(db, s.id, starPos)
			}
			return
		case err != sql.ErrNoRows:
			log.Panic(err)
		}
	}
	s.name = sysName
	err := db.QueryRow(sgx.SqlSysByName.MustSQL(), sysName).
		Scan(&s.id, null.Int64{P: &s.addr})
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
		res, err := db.Exec(sgx.SqlNewSystem.MustSQL(),
			sysName,
			null.Int64{P: &sysAddr},
			sysx, sysy, sysz,
		)
		if err != nil {
			log.Panic(err)
		}
		s.name = sysName
		countSystems++
		if id, err := res.LastInsertId(); err != nil {
			log.Panic(err)
		} else {
			s.id = int(id)
		}
		return
	case err != nil:
		log.Panic(err)
	}
	if s.addr == sqlize.NoID {
		if sysAddr != sqlize.NoID {
			// TODO sqlize
			_, err := db.Exec(`UPDATE systems SET addr=$1 WHERE id=$2`,
				sysAddr,
				s.id)
			if err != nil {
				log.Printf("failed to set system addr %d for %s: %s",
					sysAddr, sysName, err)
			}
		}
	} else if sysAddr != 0 && s.addr != sysAddr {
		log.Printf("ambg addr for %s: %d / %d", sysName, s.addr, sysAddr)
	}
}

var (
	fDB          string
	fDebug       bool
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

type Scanned struct {
	SystemAddress         int64
	StarSystem            string
	BodyID                int
	BodyName              string
	DistanceFromArrivalLS float64
}

func (scn *Scanned) LocalName() (string, bool) {
	if strings.HasPrefix(strings.ToLower(scn.BodyName), strings.ToLower(scn.StarSystem)) {
		tmp := scn.BodyName[len(scn.StarSystem):]
		return strings.TrimSpace(tmp), true
	}
	return scn.BodyName, false
}

var (
	cmdrId        int
	currentSystem system
)

var evtHdrls = map[string]func(sqlize.Tx, []byte){
	journal.FSDJumpEvent.String():   fsdJump,
	journal.DockedEvent.String():    docked,
	journal.CommanderEvent.String(): cmdrEvent,
	journal.ScanEvent.String():      scanned,
	journal.LoadGameEvent.String():  ldgEvent,
}

func importFrom(db sqlize.SQL, rd io.Reader) {
	err := sqlize.Transact(db, func(tx sqlize.Tx) error {
		scn := bufio.NewScanner(rd)
		for scn.Scan() {
			line := scn.Bytes()
			t, evt, err := events.Peek(line)
			if err != nil {
				return err
			}
			if !t.After(startAfter) {
				continue
			}
			if hdlr := evtHdrls[evt]; hdlr != nil {
				hdlr(tx, line)
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

func switchCmdr(db sqlize.Tx, cmdrNm, cmdrFid string) (cmdrId int) {
	if cmdrNm == "" {
		log.Panic("commander without name")
	}
	var fid sql.NullString
	err := db.QueryRow(sgx.SqlCmdrByName.MustSQL(), cmdrNm).
		Scan(&cmdrId, &fid)
	switch {
	case err == sql.ErrNoRows:
		res, err := db.Exec(sgx.SqlNewCmdr.MustSQL(),
			null.String{P: &cmdrFid},
			cmdrNm,
		)
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
		_, err = db.Exec(sgx.SqlCmdrSetFID.MustSQL(), cmdrFid, cmdrId)
		if err != nil {
			log.Printf("cannot update fid %s of commander %s: %s",
				cmdrFid, cmdrNm, err)
		}
	}
	return cmdrId
}

func cmdrEvent(db sqlize.Tx, line []byte) {
	var cmdr Commander
	if err := json.Unmarshal(line, &cmdr); err != nil {
		log.Panic(err)
	}
	cmdrId = switchCmdr(db, cmdr.Name, cmdr.FID)
	if cmdrId <= 0 {
		log.Panicf("no commander for jump: %s", string(line))
	}
}

func ldgEvent(db sqlize.Tx, line []byte) {
	var ldg LoadGame
	if err := json.Unmarshal(line, &ldg); err != nil {
		log.Panic(err)
	}
	currentSystem.clear()
	cmdrId = switchCmdr(db, ldg.Commander, "")
	if cmdrId <= 0 {
		log.Panicf("no commander on game load: %s", string(line))
	}
}

func fsdJump(db sqlize.Tx, line []byte) {
	if cmdrId <= 0 {
		log.Panicf("no commander for jump: %s", string(line))
	}
	var jump FSDJump
	if err := json.Unmarshal(line, &jump); err != nil {
		log.Panic(err)
	}
	if !jump.Timestamp.After(startAfter) {
		return
	}
	currentSystem.load(db, jump.StarSystem, jump.SystemAddress, &jump.StarPos)
	if currentSystem.id <= 0 {
		log.Panicf("no system for jump: %+v", &jump)
	}
	_, err := db.Exec(sgx.SqlAddVisit.MustSQL(),
		cmdrId, currentSystem.id, jump.Timestamp)
	if err != nil {
		log.Panic(err)
	}
	countVisits++
}

func docked(db sqlize.Tx, line []byte) {
	var dock Docked
	if err := json.Unmarshal(line, &dock); err != nil {
		log.Panic(err)
	}
	if !dock.Timestamp.After(startAfter) {
		return
	}
	if dock.StationName == "" {
		log.Println("no port name in ", string(line))
		return
	}
	currentSystem.load(db, dock.StarSystem, 0, nil)
	pid := getPort(db, currentSystem.id, dock.StationName, dock.StationType)
	_, err := db.Exec(sgx.SqlAddDocked.MustSQL(),
		cmdrId, pid, dock.Timestamp)
	if err != nil {
		log.Panic(err)
	}
	countDocked++
}

func scanned(db sqlize.Tx, line []byte) {
	var scn Scanned
	if err := json.Unmarshal(line, &scn); err != nil {
		log.Panic(err)
	}
	bodyLocal, _ := scn.LocalName()
	var bid int
	var typ string
	var distfa float64
	err := db.QueryRow(sgx.SqlBodyByName.MustSQL(), currentSystem.id, bodyLocal).Scan(
		&bid,
		null.String{P: &typ},
		null.Float64{P: &distfa},
	)
	if err == sql.ErrNoRows {
		_, err = db.Exec(sgx.SqlInsertBody.MustSQL(),
			currentSystem.id,
			scn.BodyID,
			bodyLocal,
			sql.NullString{},
			scn.DistanceFromArrivalLS)
		if err != nil {
			log.Panic(err)
		}
		return
	} else if err != nil {
		log.Panic(err)
	}
}

func getPort(db sqlize.Tx, sys int, name, typ string) (pid int) {
	if sys == 0 {
		log.Printf("searching port %s (%s) in system 0", name, typ)
		return 0
	}
	err := db.QueryRow(sgx.SqlPortBySysNm.MustSQL(), sys, name).Scan(&pid)
	if err == nil {
		return pid
	}
	res, err := db.Exec(sgx.SqlNewPort.MustSQL(), sys, name, strings.ToLower(typ))
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

func sysAddCoos(db sqlize.Tx, sysId int, starPos *[3]float64) {
	_, err := db.Exec(sgx.SqlSetSysCoos.MustSQL(),
		starPos[0], starPos[1], starPos[2],
		sysId)
	if err != nil {
		log.Printf("failed to set system coos for %d", sysId)
	}
}

func importLog(db sqlize.SQL, file string) {
	rd, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer rd.Close()
	importFrom(db, rd)
}

func importGz(db sqlize.SQL, file string) {
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
	flag.BoolVar(&fDebug, "debug", false, "debug import")
	flag.Parse()
	sigs := make(chan os.Signal, 1) // '1' is important for select to not always default
	signal.Notify(sigs, os.Interrupt)
	db, err := sql.Open("sqlite3", fDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var sdb sqlize.SQL
	if fDebug {
		sdb = sqlize.StdLogConnection(db)
		log.SetFlags(log.Lshortfile)
	} else {
		sdb = db
	}
	startAfter, err = lastJump(db)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
ARGS_LOOP:
	for _, arg := range flag.Args() {
		log.Printf("import %s", arg)
		switch filepath.Ext(arg) {
		case ".log":
			importLog(sdb, arg)
		case ".gz":
			importGz(sdb, arg)
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
