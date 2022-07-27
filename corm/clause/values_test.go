package clause

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	var (
		in1    = []any{[]any{"testUser1", 18}, []any{"testUser2", 20}}
		want11 = "VALUES(?,?),(?,?)"
		want12 = []any{"testUser1", 18, "testUser3", 20}
		in2    = []any{"testUser3", 52}
		want21 = "VALUES(?,?)"
		want22 = []any{"testUser3", 52}
	)
	v1, args1 := _values(in1...)
	v2, args2 := _values(in2)
	//[testUser3 52]
	//[testUser3 52]
	if v1 != want11 {
		t.Fatalf("got : %v but want : %v", v1, want11)
	} else if v2 != want21 {
		t.Fatalf("got : %v but want : %v", v2, want21)
	} else if reflect.DeepEqual(args1, want12) {
		t.Fatalf("got : %v but want : %v", args1, want12)
	} else {
		for k := range args2 {
			if args2[k] != want22[k] {
				t.Fatalf("got : %v but want : %v", args2, want22)
			}
		}
	}
}
