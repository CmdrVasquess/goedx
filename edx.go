package goedx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/CmdrVasquess/goedx/events"
	"github.com/CmdrVasquess/watched"
)

//go:generate versioner -pkg goedx -bno build_no VERSION version.go

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type Extension struct {
	JournalDir   string
	JournalAfter time.Time
	EdState      *EDState
	Galaxy       Galaxy
	CmdrFile     func(*Commander) string
	watch        *watched.JournalDir
	apps         []App
	appNms       []string
	appToks      []interface{}
}

func New(edState *EDState, gxy Galaxy) *Extension {
	if gxy == nil {
		gxy = EchoGalaxy
	}
	return &Extension{EdState: edState, Galaxy: gxy}
}

func (ext *Extension) AddApp(name string, app App) {
	ext.apps = append(ext.apps, app)
	ext.appNms = append(ext.appNms, name)
	ext.appToks = append(ext.appToks, nil)
}

func (ext *Extension) Run(latestJournal bool) (err error) {
	ext.watch = &watched.JournalDir{
		Dir:       ext.JournalDir,
		PerJLine:  ext.jLineHandler,
		OnStatChg: ext.statChangeHandler,
		Quit:      make(chan bool),
	}
	if ext.watch.Dir == "" {
		ext.watch.Dir, err = FindJournals()
		if err != nil {
			return err
		}
	}
	var latest string
	if latestJournal {
		latest, err = watched.NewestJournal(ext.watch.Dir)
		if err != nil {
			return err
		}
	}
	ext.watch.Watch(latest)
	return nil
}

func (ext *Extension) MustRun(latestJournal bool) {
	if err := ext.Run(latestJournal); err != nil {
		panic(err)
	}
}

func (ext *Extension) Stop() {
	ext.watch.Quit <- true
	<-ext.watch.Quit
}

func (ext *Extension) DiffEvtsHdls() (es []string, hs []string) {
	eventNames := events.EventNames()
	for _, enm := range eventNames {
		if _, ok := stdEvtHdlrs[enm]; !ok {
			es = append(es, enm)
		}
	}
	for hnm := range stdEvtHdlrs {
		if events.EventType(hnm) == nil {
			hs = append(hs, hnm)
		}
	}
	return es, hs
}

func (ext *Extension) jLineHandler(line []byte) {
	t, evtName, err := events.Peek(line)
	if err != nil {
		log.Errore(err)
		return
	}
	// TODO this may drop unseen events with same timestamp (1s resolution)
	if !t.After(ext.JournalAfter) {
		log.Tracea("`event` `at` outdated", evtName, t)
		return
	}
	ext.EdState.LastJournalEvent = t
	evtType := events.EventType(evtName)
	if evtType == nil {
		log.Debuga("unknown `event type`", evtName)
		return
	}
	ext.EventHandler(evtType, line)
}

func (etx *Extension) statChangeHandler(evtName string, file string) {
	evtType := events.EventType(evtName)
	if evtType == nil {
		log.Debuga("unknown `event type`", evtName)
		return
	}
	line, err := ioutil.ReadFile(file)
	if err != nil {
		log.Errora("reading `event` from `stat file`: `err`",
			evtName,
			file,
			err)
		return
	}
	etx.EventHandler(evtType, line)
}

var stdEvtHdlrs = make(map[string]func(*Extension, events.Event) Change)

type EventHandlingError struct {
	Type  events.Type
	Err   error
	Event string
}

func (ehe *EventHandlingError) Error() string {
	return fmt.Sprintf("%s '%s': [%s]", ehe.Type, ehe.Err, ehe.Event)
}

func (ext *Extension) EventHandler(evtType events.Type, raw []byte) (err error) {
	hdlr := stdEvtHdlrs[evtType.String()]
	if hdlr == nil {
		log.Debuga("no handler for `event type`", evtType)
		return err
	}
	event := evtType.New()
	if err := json.Unmarshal(raw, event); err != nil {
		log.Errora("cannot parse `event type`: `err`", evtType, err)
	}
	defer func() {
		if p := recover(); p != nil {
			evt := string(raw)
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("panic: %+v", p)
			}
			err = &EventHandlingError{
				Type:  evtType,
				Err:   err,
				Event: evt,
			}
			log.Errora("`event type` handler `panic` on `event`",
				evtType, p, evt)
		}
	}()
	for i, app := range ext.apps {
		ext.appToks[i] = ext.prepareApp(app, ext.appNms[i], event)
	}
	chg := hdlr(ext, event)
	for i, app := range ext.apps {
		if tok := ext.appToks[i]; tok != nil {
			ext.finishApp(app, ext.appNms[i], tok, event, chg)
		}
	}
	return nil
}

func (ext *Extension) prepareApp(app App, nm string, e events.Event) interface{} {
	defer func() {
		if p := recover(); p != nil {
			log.Errorf("app '%s' panics in prepare: %s", nm, p)
		}
	}()
	return app.Prepare(e)
}

func (ext *Extension) finishApp(app App, nm string, tok interface{}, e events.Event, chg Change) {
	defer func() {
		if p := recover(); p != nil {
			log.Errorf("app '%s' panic in finish: %s", nm, p)
		}
	}()
	app.Finish(tok, e, chg)
}