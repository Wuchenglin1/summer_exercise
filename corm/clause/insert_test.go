package clause

import (
	"testing"
)

func TestInsert(t *testing.T) {
	var (
		insert = []any{"testTableName", []string{"name", "age"}}
		want   = "INSERT INTO testTableName (name,age)"
		get    string
	)
	get, _ = _insert(insert...)
	if get != want {
		t.Fatalf("not equal sql ,want: %v but got : %v", want, get)
	}
}
