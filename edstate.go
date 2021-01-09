package goedx

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/CmdrVasquess/goedx/att"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/watched"
)

type HandlerFunc func(*EDState, events.Event) (att.Change, error)

var evtHdlrs = make(map[string]HandlerFunc)

const (
	ChgGame att.Change = (1 << iota)
	ChgCommander
	ChgLocation
	ChgSystem
	ChgShip
)

func SaveJSON(file string, data interface{}, logTmpl string) error {
	if !strings.HasSuffix(file, ".json") {
		file = file + ".json"
	}
	if logTmpl != "" {
		log.Infoa(logTmpl, file)
	}
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

func LoadJSON(file string, allowEmpty bool, into interface{}, logTmpl string) error {
	if !strings.HasSuffix(file, ".json") {
		file = file + ".json"
	}
	if logTmpl != "" {
		log.Infoa(logTmpl, file)
	}
	rd, err := os.Open(file)
	switch {
	case allowEmpty && os.IsNotExist(err):
		log.Warna("`file` not exists, skip loading", file)
		return nil
	case err != nil:
		return err
	}
	defer rd.Close()
	dec := json.NewDecoder(rd)
	return dec.Decode(into)
}

type Config struct {
	Galaxy          Galaxy
	CmdrFile        func(fid, name string) string
	ShutdownLogsOut bool
}

type EDState struct {
	Config

	GoEDXversion struct{ Major, Minor, Patch int }
	// Is modified w/o using Lock!
	EDVersion string
	Beta      bool
	Language  string
	L10n      struct {
		Lang   string
		Region string
	}
	Cmdr     *Commander `json:"-"`
	Loc      JSONLocation
	Ships    map[int]*Ship `json:"-"`
	JumpHist JumpHist      `json:"-"`

	lock sync.RWMutex
}

const msgNoCmdr = "no current commander"

func NewEDState() *EDState {
	res := &EDState{}
	res.GoEDXversion.Major = Major
	res.GoEDXversion.Minor = Minor
	res.GoEDXversion.Patch = Patch
	return res
}

func (ed *EDState) ResetCmdr() {
	ed.Cmdr = nil
	ed.Loc.Location = nil
	ed.Ships = make(map[int]*Ship)
	ed.JumpHist = JumpHist{
		Jumps: nil,
		Last:  0,
	}
}

func (es *EDState) SetEDVersion(v string) {
	es.EDVersion = v
	es.Beta = strings.Index(strings.ToLower(v), "beta") >= 0
}

var langMap = map[string]string{
	"English": "en",
	"German":  "de",
	"French":  "fr",
}

func ParseEDLang(edlang string) (lang, region string) {
	split := strings.Split(edlang, "\\")
	if len(split) != 2 {
		log.Errora("cannot parse `language`", edlang)
		return "", ""
	}
	lang = langMap[split[0]]
	if lang == "" {
		log.Warna("unknown `language`", split[0])
		return "", ""
	}
	return lang, split[1]
}

func (es *EDState) SetLanguage(lang string) {
	es.Language = lang
	es.L10n.Lang, es.L10n.Region = ParseEDLang(lang)
}

func (ed *EDState) SwitchCommander(fid string, name string) {
	if ed.Cmdr != nil {
		file := ed.CmdrFile(fid, name)
		if err := ed.Cmdr.Save(file); err != nil {
			log.Errore(err)
		}
	}
	ed.ResetCmdr()
	if fid == "" {
		return
	}
	if name == "" {
		log.Errora("Empty commander name for `FID`", fid)
		return
	}
	ed.Cmdr = new(Commander)
	err := LoadJSON(ed.CmdrFile(fid, name), true, ed.Cmdr, "load commander from `file`")
	if err != nil {
		log.Errore(err)
		return
	}
	ed.Cmdr.Name.Set(name, 0)
	if ed.Cmdr.FID == "" {
		ed.Cmdr.FID = fid
		return
	}
	ed.Cmdr.FID = fid
	// TODO Load ships and jump-hist
}

func (ed *EDState) MustCommander() *Commander {
	if ed.Cmdr == nil {
		panic(msgNoCmdr)
	}
	return ed.Cmdr
}

func (es *EDState) RdLocked(do func() error) error {
	es.lock.RLock()
	defer es.lock.RUnlock()
	return do()
}

func (es *EDState) WrLocked(do func() error) error {
	es.lock.Lock()
	defer es.lock.Unlock()
	return do()
}

func (ed *EDState) Save(file string, cmdrFile string) error {
	err := SaveJSON(file, ed, "save state to `file`")
	if cmdrFile != "" && ed.Cmdr != nil && ed.Cmdr.FID != "" {
		if err := ed.Cmdr.Save(cmdrFile); err != nil {
			log.Errore(err)
		}
	}
	return err
}

func (ed *EDState) Load(file string) error {
	return LoadJSON(file, true, ed, "load state from `file`")
}

func (ed *EDState) FindShip(id int) *Ship {
	if id <= 0 || id >= len(ed.Ships) {
		return nil
	}
	return ed.Ships[id]
}

func (ed *EDState) GetShip(id int) *Ship {
	res := ed.FindShip(id)
	if res == nil {
		res = new(Ship)
		ed.Ships[id] = res
	}
	return res
}

func (ed *EDState) Journal(e watched.JounalEvent) error {
	event, err := e.Event.PeekEvent()
	if err != nil {
		return err
	}
	etype := events.EventType(event)
	if etype == nil {
		log.Debuga("unknown journal `event`", event)
		return nil
	}
	eh := evtHdlrs[event]
	if eh == nil {
		log.Debuga("no handler for `event`", event)
		return nil
	}
	evt := etype.New()
	if err = json.Unmarshal(e.Event, evt); err != nil {
		return err
	}
	_, err = eh(ed, evt)
	return err
}

func (ed *EDState) Status(e watched.StatusEvent) error {
	etype := events.EventType(e.Type.String())
	if etype == nil {
		log.Debuga("unknown status `event`", etype)
		return nil
	}
	// TODO status event
	return errors.New("NYI: EDState status event")
}

func (ed *EDState) Close() error {
	// TODO status event
	return errors.New("NYI: EDState close")
}
