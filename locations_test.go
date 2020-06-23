package goedx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func ExamplePort() {
	p := Port{
		Sys: &System{
			Addr: 4711,
			Name: "Köln",
			Coos: ToSysCoos(3, 2, 1),
		},
		Name: "Hafen",
	}
	var sb bytes.Buffer
	enc := json.NewEncoder(&sb)
	enc.Encode(&p)
	os.Stdout.Write(sb.Bytes())
	sb.Reset()
	enc.Encode(JSONLocation{&p})
	os.Stdout.Write(sb.Bytes())
	var jloc JSONLocation
	fmt.Println(json.Unmarshal(sb.Bytes(), &jloc))
	sb.Reset()
	enc.Encode(&p)
	os.Stdout.Write(sb.Bytes())
	// Output:
	// {"Sys":{"Addr":4711,"Name":"Köln","Coos":[3,2,1]},"Name":"Hafen","Docked":false}
	// {"@type":"port","Docked":false,"Name":"Hafen","Sys":{"Addr":4711,"Name":"Köln","Coos":[3,2,1]}}
	// <nil>
	// {"Sys":{"Addr":4711,"Name":"Köln","Coos":[3,2,1]},"Name":"Hafen","Docked":false}
}
