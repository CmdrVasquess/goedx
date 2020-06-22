package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/CmdrVasquess/edgx"
	"github.com/CmdrVasquess/edgx/events"
)

var (
	state = edgx.NewEDState()
	ext   = edgx.New(state, edgx.EchoGalaxy)
)

func feed(rd io.Reader) {
	scn := bufio.NewScanner(rd)
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

func feedFile(name string) {
	rd, _ := os.Open(name)
	defer rd.Close()
	feed(rd)
}

func main() {
	log.SetFlags(log.Lshortfile)
	ext.CmdrFile = func(cmdr *edgx.Commander) string {
		return fmt.Sprintf("./%s.json", cmdr.FID)
	}
	state.Load("edgx-state.json")
	defer func() {
		if cmdr := ext.EdState.Cmdr; cmdr != nil {
			f := ext.CmdrFile(cmdr)
			if err := cmdr.Save(f); err != nil {
				log.Println(err)
			}
		}
		state.Save("edgx-state.json")
	}()
	if len(os.Args) < 2 {
		feed(os.Stdin)
	} else {
		for _, arg := range os.Args[1:] {
			feedFile(arg)
		}
	}
}
