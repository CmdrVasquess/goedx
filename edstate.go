package goedx

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ChgGame Change = (1 << iota)
	ChgCommander
	ChgLocation

	ChgTopNum = 3
)

func SaveJSON(file string, data interface{}) error {
	tmp := file + "~"
	wr, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer wr.Close()
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	if err = enc.Encode(data); err != nil {
		return err
	}
	wr.Close()
	return os.Rename(tmp, file)
}

func LoadJSON(file string, into interface{}) error {
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	dec := json.NewDecoder(rd)
	return dec.Decode(into)
}

type EDState struct {
	Lock sync.RWMutex `json:"-"`
	// Is modified w/o using Lock!
	EDVersion string
	Beta      bool
	Language  string
	L10n      struct {
		Lang   string
		Region string
	}
	Cmdr      *Commander `json:"-"`
	LastEvent time.Time
}

const msgNoCmdr = "no current commander"

func NewEDState() *EDState {
	res := &EDState{}
	return res
}

func (es *EDState) SetEDVersion(v string) {
	es.EDVersion = v
	es.Beta = strings.Index(strings.ToLower(v), "beta") >= 0
}

var langMap = map[string]string{
	"English": "en",
}

func (es *EDState) SetLanguage(lang string) {
	es.Language = lang
	split := strings.Split(lang, "\\")
	if len(split) != 2 {
		log.Errora("cannot partse `language`", lang)
		es.L10n.Lang = ""
		es.L10n.Region = ""
	}
	es.L10n.Lang = langMap[split[0]]
	if es.L10n.Lang == "" {
		log.Warna("unknown `language`", split[0])
	}
	es.L10n.Region = split[1]
}

func (es *EDState) MustCommander() *Commander {
	if es.Cmdr == nil {
		panic(msgNoCmdr)
	}
	return es.Cmdr
}

func (es *EDState) Read(do func() error) error {
	es.Lock.RLock()
	defer es.Lock.RUnlock()
	return do()
}

func (es *EDState) Write(do func() error) error {
	es.Lock.Lock()
	defer es.Lock.Unlock()
	return do()
}

func (es *EDState) WriteCmdr(do func(*Commander) error) error {
	es.Lock.Lock()
	defer es.Lock.Unlock()
	if es.Cmdr == nil {
		return errors.New(msgNoCmdr)
	}
	return do(es.Cmdr)
}

func (ed *EDState) Save(file string) error {
	log.Infoa("save state to `file`", file)
	return SaveJSON(file, ed)
}

func (ed *EDState) Load(file string) error {
	log.Infoa("load state from `file`", file)
	return LoadJSON(file, ed)
}

type Commander struct {
	FID    string
	Name   string
	Ranks  Ranks
	ShipID int
	At     JSONLocation
	Ships  map[int]*Ship
	Mats   Materials
	inShip *Ship
}

func NewCommander() *Commander {
	return &Commander{
		Ships: make(map[int]*Ship),
	}
}

func (cmdr *Commander) FindShip(id int) *Ship {
	if id <= 0 {
		return nil
	}
	return cmdr.Ships[id]
}

func (cmdr *Commander) GetShip(id int) *Ship {
	res := cmdr.FindShip(id)
	if res == nil {
		res = new(Ship)
		cmdr.Ships[id] = res
	}
	return res
}

func (cmdr *Commander) SetShip(id int) *Ship {
	if id < 0 {
		cmdr.inShip = nil
		cmdr.ShipID = -1
		return nil
	}
	res := cmdr.GetShip(id)
	cmdr.ShipID = id
	cmdr.inShip = res
	res.Berth = nil
	return res
}

func (cmdr *Commander) StoreCurrentShip() {
	cmdr.ShipID = -1
	if cmdr.inShip == nil {
		return
	}
	ship := cmdr.inShip
	cmdr.inShip = nil
	if port := cmdr.At.Port(); port != nil {
		ship.Berth = port
	}
}

func (cmdr *Commander) Save(file string) error {
	log.Infoa("save `commander` with `fid` to `file`", cmdr.Name, cmdr.FID, file)
	return SaveJSON(file, cmdr)
}

func (cmdr *Commander) Load(file string) error {
	log.Infoa("load commander from `file`", file)
	err := LoadJSON(file, cmdr)
	cmdr.inShip = cmdr.FindShip(cmdr.ShipID)
	return err
}

type Rank struct {
	Level    int
	Progress int
}

//go:generate stringer -type RankType
type RankType int

const (
	Combat RankType = iota
	Trade
	Explore
	CQC
	Federation
	Empire

	RanksNum
)

type Ranks [RanksNum]Rank

type Ship struct {
	Type  string
	Ident string
	Name  string
	Berth *Port      `json:",omitempty"`
	Sold  *time.Time `json:",omitempty"`
}

type Materials struct {
	Raw map[string]Material
	Man map[string]Material
	Enc map[string]Material
}

type Material struct {
	Stock  int
	Demand int
}
