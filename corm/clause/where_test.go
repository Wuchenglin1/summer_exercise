package clause

import (
	"testing"
)

func TestWhere(t *testing.T) {
	var (
		in       = []any{[]any{"Name = ?", "testUser"}}
		wantSql  = "WHERE Name = ?"
		wantArgs = []any{"testUser"}
	)
	gotSql, gotArgs := _where(in...)
	if gotSql != wantSql {
		t.Fatalf("gotSql %v but wantSql %v", gotSql, wantSql)
	} else {
		if len(gotArgs) != len(wantArgs) {
			t.Fatalf("gotArgs %v but wantArgs %v", gotArgs, wantArgs)
		}
		for k := range gotArgs {
			if gotArgs[k] != wantArgs[k] {
				t.Fatalf("gotArgs %v but wantArgs %v", gotArgs, wantArgs)
			}
		}

	}
}
