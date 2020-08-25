package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	t := template.Must(template.New("journal").Parse(tmpl))
	flag.Parse()
	var args struct {
		Capitals string
		Downcase string
	}
	for _, arg := range flag.Args() {
		args.Capitals = arg
		args.Downcase = strings.ToLower(arg)
		fnm := args.Downcase + ".go"
		if _, err := os.Stat(fnm); !os.IsNotExist(err) {
			log.Fatalf("file '%s' already exists", fnm)
		}
		wr, err := os.Create(fnm)
		if err != nil {
			log.Fatal(err)
		}
		t.Execute(wr, &args)
		wr.Close()
	}
}

const tmpl = `package journal

import "github.com/CmdrVasquess/goedx/events"

type {{.Downcase}}T string

const {{.Capitals}}Event = {{.Downcase}}T("{{.Capitals}}")

func (t {{.Downcase}}T) New() events.Event { return new({{.Capitals}}) }
func (t {{.Downcase}}T) String() string    { return string(t) }

type {{.Capitals}} struct {
	events.Common
	// TODO fill in event attributes
}

func (_ *{{.Capitals}}) EventType() events.Type { return {{.Capitals}}Event }

func init() {
	events.RegisterType(string({{.Capitals}}Event), {{.Capitals}}Event)
}
`
