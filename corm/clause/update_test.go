package clause

import (
	"testing"
)

func TestUpdate(t *testing.T) {
	var (
		uMap     = make(map[string]any)
		wantSql  = "UPDATE tableName SET Name= ?, Age= ?"
		wantArgs = []any{"hello", 18}
	)
	uMap["Name"] = "hello"
	uMap["Age"] = 18
	var in = []any{"tableName", uMap}
	got, args := _update(in...)
	if got != wantSql {
		t.Fatalf("gotSql : %v but wantSql : %v", got, wantSql)
	} else {
		if len(args) != len(wantArgs) {
			t.Fatalf("gotArgs : %v but wantArgs : %v", args, wantArgs)
		}
		for k := range args {
			if args[k] != wantArgs[k] {
				t.Fatalf("gotArgs : %v but wantArgs : %v", args, wantArgs)
			}
		}
	}
}
