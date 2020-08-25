package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"git.fractalqb.de/fractalqb/c4hgol"

	"git.fractalqb.de/fractalqb/qbsllm"
	"github.com/CmdrVasquess/goedx"
	"github.com/CmdrVasquess/goedx/apps/bboltgalaxy"
	"github.com/CmdrVasquess/goedx/apps/l10n"
	"github.com/CmdrVasquess/goedx/events"
)

var (
	log    = qbsllm.New(qbsllm.Linfo, "feed", nil, nil)
	logCgf = c4hgol.Config(qbsllm.NewConfig(log), goedx.LogCfg)

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
			log.Errore(err)
			log.Fatalf("[%s]", string(line))
		} else if et := events.EventType(evt); et != nil {
			if err = ext.EventHandler(et, line); err != nil {
				log.Fatale(err)
			}
		}
	}
}

func feedFile(name string) {
	rd, _ := os.Open(name)
	defer rd.Close()
	feed(rd)
}

func main() {
	c4hgol.ShowSource(logCgf, true, false)
	gxyFile := flag.String("galaxy", "", "set galaxy database file")
	l10nDir := flag.String("l10n", "", "load/save l10ns to dir")
	flag.Parse()
	if *gxyFile != "" {
		gxy, err := bboltgalaxy.Open(*gxyFile)
		if err == nil {
			ext.Galaxy = gxy
			ext.AddApp("bbgxy", gxy)
			defer gxy.Close()
		}
	}
	if *l10nDir != "" {
		l10nApp := l10n.New(*l10nDir, state)
		ext.AddApp("l10n", l10nApp)
		defer l10nApp.Close()
	}
	{
		e, h := ext.DiffEvtsHdls()
		log.Infof("events w/o handler: %s", e)
		log.Infof("handlers w/o event: %s", h)
	}
	ext.CmdrFile = func(cmdr *goedx.Commander) string {
		return fmt.Sprintf("./%s.json", cmdr.FID)
	}
	state.Load("goedx-state.json")
	defer func() {
		if cmdr := ext.EDState.Cmdr; cmdr != nil && cmdr.FID != "" {
			f := ext.CmdrFile(cmdr)
			if err := cmdr.Save(f); err != nil {
				log.Errore(err)
			}
		}
		state.Save("goedx-state.json", "")
	}()
	if len(os.Args) < 2 {
		feed(os.Stdin)
	} else {
		for _, arg := range flag.Args() {
			feedFile(arg)
		}
	}
}
