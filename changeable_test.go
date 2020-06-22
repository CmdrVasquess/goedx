package edgx

import (
	"encoding/json"
	"math"
	"testing"
)

func TestChgF32_json(t *testing.T) {
	test := func(f float32) {
		cf := ChgF32(f)
		data, err := json.Marshal(cf)
		if err != nil {
			t.Error(err)
		}
		var rf ChgF32
		if err = json.Unmarshal(data, &rf); err != nil {
			t.Error(err)
		}
		if math.IsNaN(float64(rf)) && math.IsNaN(float64(f)) {
			return
		}
		if float32(rf) != f {
			t.Errorf("expect %f, got %f", f, rf)
		}
	}
	test(4711)
	test(float32(math.Inf(-1)))
	test(float32(math.Inf(1)))
	test(float32(math.NaN()))
}
