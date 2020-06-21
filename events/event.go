package events

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

var (
	timestampTag = []byte(`"timestamp":`)
	eventTag     = []byte(`"event":`)
)

func Peek(str []byte) (t time.Time, event string, err error) {
	idx := bytes.Index(str, timestampTag)
	if idx < 0 {
		return time.Time{}, "", errors.New("no timestamp in event")
	}
	val := str[idx+13 : idx+33]
	t, err = time.Parse(time.RFC3339, string(val))
	if err != nil {
		panic(err)
	}
	str = str[idx+35:]
	idx = bytes.Index(str, eventTag)
	if idx < 0 {
		return time.Time{}, "", errors.New("no event type in event")
	}
	str = str[idx+9:]
	idx = bytes.IndexByte(str, '"')
	if idx < 0 {
		panic("cannot find end of event type")
	}
	return t, string(str[:idx]), nil
}

type Event interface {
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

var eventTypes = make(map[string]Type)
