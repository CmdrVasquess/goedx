package goedx

import (
	"encoding/json"
	"math"
	"strconv"
)

type Change uint64

func (chg Change) Any(c Change) bool       { return (chg & c) != 0 }
func (chg Change) All(c Change) bool       { return (chg & c) == c }
func (chg Change) Without(c Change) Change { return chg &^ c }

type ChgString string

func (s *ChgString) Set(v string, chg Change) Change {
	if string(*s) != v {
		*s = ChgString(v)
		return chg
	}
	return 0
}

type ChgInt int

func (i *ChgInt) Set(v int, chg Change) Change {
	if int(*i) != v {
		*i = ChgInt(v)
		return chg
	}
	return 0
}

type ChgF32 float32

func (f *ChgF32) Set(v float32, chg Change) Change {
	if float32(*f) != v {
		*f = ChgF32(v)
		return chg
	}
	return 0
}

func (f ChgF32) MarshalJSON() ([]byte, error) {
	x := float64(f)
	switch {
	case math.IsNaN(x):
		return json.Marshal("NaN")
	case math.IsInf(x, 1):
		return json.Marshal("+∞")
	case math.IsInf(x, -1):
		return json.Marshal("-∞")
	default:
		return strconv.AppendFloat(nil, x, 'f', -1, 32), nil
	}
}

func (f *ChgF32) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"NaN"`:
		*f = ChgF32(math.NaN())
	case `"+∞"`:
		*f = ChgF32(math.Inf(1))
	case `"-∞"`:
		*f = ChgF32(math.Inf(-1))
	default:
		x, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}
		*f = ChgF32(x)
	}
	return nil
}

type ChgF64 float64

func (f *ChgF64) Set(v float64, chg Change) Change {
	if float64(*f) != v {
		*f = ChgF64(v)
		return chg
	}
	return 0
}

func (f ChgF64) MarshalJSON() ([]byte, error) {
	x := float64(f)
	switch {
	case math.IsNaN(x):
		return json.Marshal("Nan")
	case math.IsInf(x, 1):
		return json.Marshal("+∞")
	case math.IsInf(x, -1):
		return json.Marshal("-∞")
	default:
		return strconv.AppendFloat(nil, x, 'f', -1, 64), nil
	}
}

func (f *ChgF64) UnmarshalJSON(data []byte) error {
	str := string(data)
	switch str {
	case `"NaN"`:
		*f = ChgF64(math.NaN())
	case `"+∞"`:
		*f = ChgF64(math.Inf(1))
	case `"-∞"`:
		*f = ChgF64(math.Inf(-1))
	default:
		x, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		*f = ChgF64(x)
	}
	return nil
}

type SysCoos [3]ChgF32

func ToSysCoos(x, y, z float32) SysCoos {
	return SysCoos{ChgF32(x), ChgF32(y), ChgF32(z)}
}

func (sc *SysCoos) Set(x, y, z float32, chg Change) (res Change) {
	res |= sc[0].Set(x, chg)
	res |= sc[1].Set(y, chg)
	res |= sc[2].Set(z, chg)
	return chg
}

func (sc *SysCoos) Valid() bool {
	return !math.IsNaN(float64((*sc)[0])) &&
		!math.IsNaN(float64((*sc)[1])) &&
		!math.IsNaN(float64((*sc)[2]))
}
