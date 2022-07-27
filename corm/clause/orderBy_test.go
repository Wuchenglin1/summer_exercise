package clause

import "testing"

func TestOrderBy(t *testing.T) {
	var (
		in   = "name"
		want = "ORDER BY name"
	)
	if got, _ := _orderBy(in); got != want {
		t.Fatalf("got : %v but want : %v", got, want)
	}
}
