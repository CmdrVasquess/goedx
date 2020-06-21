package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/CmdrVasquess/edgx"
	"github.com/CmdrVasquess/edgx/events"
)

var (
	state = edgx.NewEDState()
	ext   = edgx.New(state)
)

func main() {
	ext.CmdrFile = func(cmdr *edgx.Commander) string {
		return fmt.Sprintf("./%s.json", cmdr.FID)
	}
	state.Load("edgx-state.json")
	defer func() {
		if cmdr := ext.EdState.Cmdr; cmdr != nil {
			f := ext.CmdrFile(cmdr)
			cmdr.Save(f)
		}
		state.Save("edgx-state.json")
	}()
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		line := scn.Bytes()
		_, evt, err := events.Peek(line)
		if err != nil {
			log.Println(err)
		} else if et := events.EventType(evt); et == nil {
			log.Printf("unknown event type: '%s'", evt)
		} else {
			ext.EventHandler(et, line)
		}
	}
}
