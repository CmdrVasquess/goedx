package bboltgalaxy

import (
	"reflect"
	"testing"

	"github.com/CmdrVasquess/goedx"
)

var _ goedx.Galaxy = (*Galaxy)(nil)

func TestSystem_StoreLoad(t *testing.T) {
	gxy, err := Open(t.Name() + ".galaxy")
	if err != nil {
		t.Fatal(err)
	}
	s1 := &System{
		System: *goedx.NewSystem(4711, "4711"),
	}
	gxy.UpdateSystem(s1)
	s2 := gxy.FindSystemByAddr(4711)
	if s2 == nil {
		t.Fatal("cannot load system 4711")
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Error("systems differ")
	}
}
