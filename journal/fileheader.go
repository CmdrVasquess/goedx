package journal

import "github.com/CmdrVasquess/goedx/events"

type fileheaderT string

const FileheaderEvent = fileheaderT("Fileheader")

func (t fileheaderT) New() events.Event { return new(Fileheader) }
func (t fileheaderT) String() string    { return string(t) }

type Fileheader struct {
	events.Common
	GameVersion string `json:"gameversion"`
	Language    string `json:"language"`
}

func init() {
	events.RegisterType(string(FileheaderEvent), FileheaderEvent)
}
