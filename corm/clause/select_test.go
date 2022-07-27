package clause

import "testing"

func TestSelect(t *testing.T) {
	var (
		in   = []any{"tableName", []string{"name", "age"}}
		want = "SELECT name,age FROM tableName"
	)
	got, _ := _select(in...)
	if got != want {
		t.Logf("got : %v but want : %v", got, want)
	}
}
