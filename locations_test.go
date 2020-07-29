package goedx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func ExampleSystem() {
	s := System{
		Addr: 4711,
		Name: "Köln",
		Coos: ToSysCoos(3, 2, 1),
	}
	var sb bytes.Buffer
	enc := json.NewEncoder(&sb)
	enc.SetIndent("", "  ")
	enc.Encode(&s)
	os.Stdout.Write(sb.Bytes())
	sb.Reset()
	enc.Encode(JSONLocation{&s})
	os.Stdout.Write(sb.Bytes())
	var jloc JSONLocation
	fmt.Println(json.Unmarshal(sb.Bytes(), &jloc))
	sb.Reset()
	enc.Encode(&s)
	os.Stdout.Write(sb.Bytes())
	// Output:
	// {
	//   "Addr": 4711,
	//   "Name": "Köln",
	//   "Coos": [
	//     3,
	//     2,
	//     1
	//   ],
	//   "FirstAccess": "0001-01-01T00:00:00Z",
	//   "LastAccess": "0001-01-01T00:00:00Z"
	// }
	// {
	//   "@type": "system",
	//   "Addr": 4711,
	//   "Coos": [
	//     3,
	//     2,
	//     1
	//   ],
	//   "Name": "Köln"
	// }
	// <nil>
	// {
	//   "Addr": 4711,
	//   "Name": "Köln",
	//   "Coos": [
	//     3,
	//     2,
	//     1
	//   ],
	//   "FirstAccess": "0001-01-01T00:00:00Z",
	//   "LastAccess": "0001-01-01T00:00:00Z"
	// }
}

func ExamplePort() {
	p := Port{
		Sys: &System{
			Addr: 4711,
			Name: "Köln",
			Coos: ToSysCoos(3, 2, 1),
		},
		Name:   "Hafen",
		Type:   "Orbis",
		Docked: true,
	}
	var sb bytes.Buffer
	enc := json.NewEncoder(&sb)
	enc.SetIndent("", "  ")
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
	// {
	//   "Sys": {
	//     "Addr": 4711,
	//     "Name": "Köln",
	//     "Coos": [
	//       3,
	//       2,
	//       1
	//     ],
	//     "FirstAccess": "0001-01-01T00:00:00Z",
	//     "LastAccess": "0001-01-01T00:00:00Z"
	//   },
	//   "Name": "Hafen",
	//   "Type": "Orbis",
	//   "Docked": true
	// }
	// {
	//   "@type": "port",
	//   "Docked": true,
	//   "Name": "Hafen",
	//   "Sys": {
	//     "Addr": 4711,
	//     "Coos": [
	//       3,
	//       2,
	//       1
	//     ],
	//     "Name": "Köln"
	//   },
	//   "Type": "Orbis"
	// }
	// <nil>
	// {
	//   "Sys": {
	//     "Addr": 4711,
	//     "Name": "Köln",
	//     "Coos": [
	//       3,
	//       2,
	//       1
	//     ],
	//     "FirstAccess": "0001-01-01T00:00:00Z",
	//     "LastAccess": "0001-01-01T00:00:00Z"
	//   },
	//   "Name": "Hafen",
	//   "Type": "Orbis",
	//   "Docked": true
	// }
}
