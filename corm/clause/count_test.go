package clause

import "testing"

func TestCount(t *testing.T) {
	var (
		in   = "tableName"
		want = "SELECT COUNT(*) FROM tableName"
	)
	got, _ := _count(in)
	if got != want {
		t.Fatalf("got %v but want %v", got, want)
	}
}
