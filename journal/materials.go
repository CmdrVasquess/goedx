package journal

import "github.com/CmdrVasquess/edgx/events"

type materialsT string

const MaterialsEvent = materialsT("Materials")

func (t materialsT) New() events.Event { return new(Materials) }
func (t materialsT) String() string    { return string(t) }

type Materials struct {
	events.Common
	Raw          []Material
	Manufactured []Material
	Encoded      []Material
}

type Material struct {
	Name    string
	NameL7d string `json:"Name_Localised,omitempty"`
	Count   int
}

func init() {
	events.RegisterType(string(MaterialsEvent), MaterialsEvent)
}
