package goedx

import "github.com/CmdrVasquess/goedx/events"

type App interface {
	PrepareEDEvent(e events.Event) (token interface{})
	FinishEDEvent(token interface{}, e events.Event, chg Change)
}

type AppChannel struct {
	app App
	c   chan finish
}

func NewAppChannel(app App, capacity int) *AppChannel {
	res := &AppChannel{
		app: app,
		c:   make(chan finish, capacity),
	}
	go res.run()
	return res
}

func (ac *AppChannel) Close() {
	close(ac.c)
}

func (ac *AppChannel) Finish(token interface{}, e events.Event, chg Change) {
	ac.c <- finish{token, e, chg}
}

type finish struct {
	tok interface{}
	evt events.Event
	chg Change
}

func (ac *AppChannel) run() {
	for f := range ac.c {
		ac.app.FinishEDEvent(f.tok, f.evt, f.chg)
	}
}
