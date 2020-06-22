package edgx

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/CmdrVasquess/edgx/events"
	"github.com/CmdrVasquess/watched"
)

type Extension struct {
	JournalDir   string
	JournalAfter time.Time
	EdState      *EDState
	Galaxy       Galaxy
	CmdrFile     func(*Commander) string
	watch        *watched.JournalDir
}

func New(edState *EDState, gxy Galaxy) *Extension {
	if gxy == nil {
		gxy = EchoGalaxy
	}
	return &Extension{EdState: edState, Galaxy: gxy}
}

func (edgx *Extension) Run(latestJournal bool) (err error) {
	edgx.watch = &watched.JournalDir{
		Dir:       edgx.JournalDir,
		PerJLine:  edgx.jLineHandler,
		OnStatChg: edgx.statChangeHandler,
	}
	if edgx.watch.Dir == "" {
		edgx.watch.Dir, err = FindJournals()
		if err != nil {
			return err
		}
	}
	var latest string
	if latestJournal {
		latest, err = watched.NewestJournal(edgx.watch.Dir)
		if err != nil {
			return err
		}
	}
	edgx.watch.Watch(latest)
	return nil
}

func (edgx *Extension) MustRun(latestJournal bool) {
	if err := edgx.Run(latestJournal); err != nil {
		panic(err)
	}
}

func (edgx *Extension) jLineHandler(line []byte) {
	t, evtName, err := events.Peek(line)
	if err != nil {
		log.Errore(err)
		return
	}
	if !t.After(edgx.JournalAfter) {
		log.Tracea("`event` `at` outdated", evtName, t)
		return
	}
	evtType := events.EventType(evtName)
	if evtType == nil {
		log.Debuga("unknown `event type`", evtName)
		return
	}
	edgx.EventHandler(evtType, line)
}

func (edgx *Extension) statChangeHandler(evtName string, file string) {
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
	edgx.EventHandler(evtType, line)
}

var stdEvtHdlrs = make(map[string]func(*Extension, events.Event) Change)

func (edgx *Extension) EventHandler(evtType events.Type, raw []byte) {
	hdlr := stdEvtHdlrs[evtType.String()]
	if hdlr == nil {
		log.Debuga("no handler for `event type`", evtType)
		return
	}
	event := evtType.New()
	if err := json.Unmarshal(raw, event); err != nil {
		log.Errora("cannot parse `event type`: `err`", evtType, err)
	}
	defer func() {
		if p := recover(); p != nil {
			log.Errora("`event type` handler `panic` on `event`",
				evtType,
				p,
				string(raw))
		}
	}()
	hdlr(edgx, event)
	edgx.EdState.LastEvent = event.Timestamp()
}
