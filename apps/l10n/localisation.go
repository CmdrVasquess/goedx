package l10n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"git.fractalqb.de/fractalqb/c4hgol"
	"git.fractalqb.de/fractalqb/qbsllm"
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/goedx/journal"
)

var (
	log    = qbsllm.New(qbsllm.Lnormal, "goedx.app.l10n", nil, nil)
	LogCfg = c4hgol.Config(qbsllm.NewConfig(log))
)

type Locales struct {
	BaseDir     string
	edState     *goedx.EDState
	currentLang string
	Shiptypes   map[string]string
	MatsRaw     map[string]string
	MatsMan     map[string]string
	MatsEnc     map[string]string
	Economies   map[string]string
	Securities  map[string]string
}

func New(dir string, edState *goedx.EDState) *Locales {
	res := &Locales{BaseDir: dir, edState: edState}
	res.clearMaps()
	return res
}

func (loc *Locales) Close() {
	loc.save()
}

func (loc *Locales) ShipType(key string) (string, bool) {
	return loc.local(loc.Shiptypes, key)
}

func (loc *Locales) Economy(key string) (string, bool) {
	return loc.local(loc.Economies, key)
}

func (loc *Locales) Security(key string) (string, bool) {
	return loc.local(loc.Securities, key)
}

func (loc *Locales) RawMaterial(key string) (string, bool) {
	return loc.local(loc.MatsRaw, key)
}

func (loc *Locales) ManMaterial(key string) (string, bool) {
	return loc.local(loc.MatsMan, key)
}

func (loc *Locales) EncMaterial(key string) (string, bool) {
	return loc.local(loc.MatsEnc, key)
}

func (loc *Locales) local(m map[string]string, key string) (string, bool) {
	res := m[key]
	if res == "" {
		return key, false
	}
	return res, true
}

func getLang(edlang string) string {
	lang, region := goedx.ParseEDLang(edlang)
	switch {
	case lang == "":
		return ""
	case region == "":
		return lang
	}
	return fmt.Sprintf("%s-%s", lang, region)
}

func normKey(key string) string {
	if key[0] == '$' {
		key = key[1:]
	}
	if l := len(key) - 1; key[l] == ';' {
		key = key[:l]
	}
	return key
}

func (loc *Locales) PrepareEDEvent(e events.Event) interface{} {
	switch evt := e.(type) {
	case *journal.Fileheader:
		lang := getLang(evt.Language)
		if lang == "" {
			return nil
		}
		return lang
	case *journal.Materials, *journal.FSDJump, *journal.ShipTargeted:
		if loc.edState.L10n.Lang == "" {
			return nil
		}
		return true
	}
	return nil
}

func (loc *Locales) FinishEDEvent(token interface{}, e events.Event, chg goedx.Change) {
	switch evt := e.(type) {
	case *journal.ShipTargeted:
		loc.finishShipTargeted(evt, chg)
	case *journal.FSDJump:
		loc.finishFSDJump(evt, chg)
	case *journal.Fileheader:
		loc.save()
		loc.load(token.(string))
	case *journal.Materials:
		loc.finishMaterials(evt, chg)
	default:
		log.Errora("drop invalid `event type`", reflect.TypeOf(e))
	}
}

const (
	mapShiptype    = "shiptype"
	mapEconomy     = "economy"
	mapSecurity    = "security"
	mapMatNamesRaw = "matnames-raw"
	mapMatNamesMan = "matnames-man"
	mapMatNamesEnc = "matnames-enc"
)

func (loc *Locales) save() {
	if loc.currentLang == "" {
		return
	}
	log.Debuga("saving current `lang` to `dir`", loc.currentLang, loc.BaseDir)
	loc.saveMap(mapShiptype, loc.Shiptypes)
	loc.saveMap(mapEconomy, loc.Economies)
	loc.saveMap(mapSecurity, loc.Securities)
	loc.saveMap(mapMatNamesRaw, loc.MatsRaw)
	loc.saveMap(mapMatNamesMan, loc.MatsMan)
	loc.saveMap(mapMatNamesEnc, loc.MatsEnc)
}

func (loc *Locales) saveMap(name string, m map[string]string) {
	file := loc.mapFile(name)
	log.Tracea("save `map` to `file`", name, file)
	tmp := file + "~"
	wr, err := os.Create(tmp)
	if err != nil {
		log.Errora("create `map` `err`", name, err)
		return
	}
	defer wr.Close()
	enc := json.NewEncoder(wr)
	enc.SetIndent("", "\t")
	if err = enc.Encode(m); err != nil {
		log.Errora("write `map` `err`", name, err)
		return
	}
	wr.Close()
	if err = os.Rename(tmp, file); err != nil {
		log.Errore(err)
	}
}

func (loc *Locales) clearMaps() {
	log.Debugs("clearing maps")
	loc.Shiptypes = make(map[string]string)
	loc.Economies = make(map[string]string)
	loc.Securities = make(map[string]string)
	loc.MatsRaw = make(map[string]string)
	loc.MatsMan = make(map[string]string)
	loc.MatsEnc = make(map[string]string)
}

func (loc *Locales) load(lang string) {
	if lang == "" {
		loc.clearMaps()
		loc.currentLang = ""
		return
	}
	log.Debuga("load `lang` from `dir`", lang, loc.BaseDir)
	loc.currentLang = lang
	loc.Shiptypes = loc.loadMap(mapShiptype)
	loc.Economies = loc.loadMap(mapEconomy)
	loc.Securities = loc.loadMap(mapSecurity)
	loc.MatsRaw = loc.loadMap(mapMatNamesRaw)
	loc.MatsMan = loc.loadMap(mapMatNamesMan)
	loc.MatsEnc = loc.loadMap(mapMatNamesEnc)
}

func (loc *Locales) loadMap(name string) map[string]string {
	res := make(map[string]string)
	file := loc.mapFile(name)
	log.Tracea("load `map` from `file`", name, file)
	rd, err := os.Open(file)
	if os.IsNotExist(err) {
		log.Warna("no `map` `file`", name, file)
		return res
	}
	if err != nil {
		log.Errore(err)
	}
	defer rd.Close()
	dec := json.NewDecoder(rd)
	if err = dec.Decode(&res); err != nil {
		log.Errore(err)
	}
	return res
}

func (loc *Locales) mapFile(name string) string {
	res := filepath.Join(loc.BaseDir, loc.currentLang)
	if _, err := os.Stat(res); os.IsNotExist(err) {
		log.Infoa("create `lang` directory", loc.currentLang)
		if err = os.Mkdir(res, 0777); err != nil {
			log.Errore(err)
		}
	}
	return filepath.Join(loc.BaseDir, loc.currentLang, name+".json")
}
