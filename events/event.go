package events

import (
	"fmt"
	"time"
)

type Event interface {
	EventType() Type
	Timestamp() time.Time
	Event() string
}

type Common struct {
	Time time.Time `json:"timestamp"`
	Tag  string    `json:"event"`
}

func (c *Common) Timestamp() time.Time { return c.Time }

func (c *Common) Event() string { return c.Tag }

type Type interface {
	fmt.Stringer
	New() Event
}

func RegisterType(name string, t Type) {
	eventTypes[name] = t
}

func EventType(event string) Type {
	return eventTypes[event]
}

func EventNames() []string {
	res := make([]string, 0, len(eventTypes))
	for nm := range eventTypes {
		res = append(res, nm)
	}
	return res
}

var eventTypes = make(map[string]Type)
