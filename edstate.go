package edgx

import (
	"encoding/json"
	"os"
	"sync"
	"time"
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
	LastEvent time.Time
	Cmdr      *Commander `json:"-"`
}

func (es *EDState) MustCommander() *Commander {
	if es.Cmdr == nil {
		panic("no current commander")
	}
	return es.Cmdr
}

func NewEDState() *EDState {
	res := &EDState{}
	return res
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
	ShipID int
	Loc    JSONLocation
	Ships  []*Ship
	Mats   Materials
	inShip *Ship
}

func (cmdr *Commander) FindShip(id int) *Ship {
	if id < 0 {
		return nil
	}
	if cmdr.ShipID == id {
		if cmdr.inShip == nil || cmdr.inShip.ID != id {
			panic("assert failed: ship id mismatch")
		}
		return cmdr.inShip
	}
	for i := range cmdr.Ships {
		s := cmdr.Ships[i]
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (cmdr *Commander) GetShip(id int) *Ship {
	res := cmdr.FindShip(id)
	if res == nil {
		res := &Ship{ID: id}
		cmdr.Ships = append(cmdr.Ships, res)
	}
	return res
}

func (cmdr *Commander) SetShip(id int) *Ship {
	if id < 0 {
		cmdr.inShip = nil
		return nil
	}
	res := cmdr.GetShip(id)
	cmdr.ShipID = res.ID
	cmdr.inShip = res
	res.Berth = nil
	return res
}

func (cmdr *Commander) Save(file string) error {
	log.Infoa("save `commander` with `fid` to `file`", cmdr.Name, cmdr.FID, file)
	return SaveJSON(file, cmdr)
}

func (cmdr *Commander) Load(file string) error {
	log.Infoa("load commander from `file`", file)
	return LoadJSON(file, cmdr)
}

type Ship struct {
	ID    int
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
