package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/events"
)

var (
	state = goedx.NewEDState()
	ext   = goedx.New(state, goedx.EchoGalaxy)
)

func feed(rd io.Reader) {
	scn := bufio.NewScanner(rd)
	for scn.Scan() {
		line := scn.Bytes()
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		_, evt, err := events.Peek(line)
		if err != nil {
			log.Print(err)
			log.Fatalf("[%s]", string(line))
		} else if et := events.EventType(evt); et == nil {
			log.Printf("unknown event type: '%s'", evt)
		} else if err = ext.EventHandler(et, line); err != nil {
			log.Fatal(err)
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
	{
		e, h := ext.DiffEvtsHdls()
		log.Println("events w/o handler:", e)
		log.Println("handlers w/o event:", h)
	}
	ext.CmdrFile = func(cmdr *goedx.Commander) string {
		return fmt.Sprintf("./%s.json", cmdr.FID)
	}
	state.Load("goedx-state.json")
	defer func() {
		if cmdr := ext.EdState.Cmdr; cmdr != nil {
			f := ext.CmdrFile(cmdr)
			if err := cmdr.Save(f); err != nil {
				log.Println(err)
			}
		}
		state.Save("goedx-state.json")
	}()
	if len(os.Args) < 2 {
		feed(os.Stdin)
	} else {
		for _, arg := range os.Args[1:] {
			feedFile(arg)
		}
	}
}
