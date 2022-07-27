package clause

import "testing"

func TestLimit(t *testing.T) {
	var (
		in       = 5
		wantSql  = "LIMIT ?"
		wantArgs = []any{5}
	)
	gotSql, gotArgs := _limit(in)
	if gotSql != wantSql {
		t.Logf("got : %v but want : %v", gotSql, wantSql)
	} else {
		if len(gotArgs) != len(wantArgs) {
			t.Logf("got : %v but want : %v", gotArgs, wantArgs)
		}
		for k := range gotArgs {
			if gotArgs[k] != wantArgs[k] {
				t.Logf("got : %v but want : %v", gotArgs, wantArgs)
			}
		}
	}
}
