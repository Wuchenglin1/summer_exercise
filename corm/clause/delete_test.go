package clause

import "testing"

func TestDelete(t *testing.T) {
	var (
		in   = "tableName"
		want = "DELETE FROM tableName "
	)
	sql, _ := _delete(in)
	if sql != want {
		t.Fatalf("got : %v but got : %v", sql, want)
	}
}
