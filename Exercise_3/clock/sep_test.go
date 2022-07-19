package clock

import (
	"fmt"
	"testing"
	"time"
)

func TestSepSchedule_Next(t *testing.T) {
	type fields struct {
		Second []Schedule
		Minute []Schedule
		Hour   []Schedule
		Day    []Schedule
		Month  []Schedule
		Dow    []Schedule
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{"test", fields{
			[]Schedule{{1, 59, true, 10}},
			[]Schedule{{0, 59, true, 2}},
			[]Schedule{{0, 59, false, 1}},
			[]Schedule{{1, 31, false, 1}},
			[]Schedule{{1, 12, false, 1}},
			[]Schedule{{1, 7, false, 1}},
		}, args{time.Now()}, time.Now().Add(time.Second * 5)},
	}
	for _, tt := range tests {
		s := &SepSchedule{
			Second: tt.fields.Second,
			Minute: tt.fields.Minute,
			Hour:   tt.fields.Hour,
			Day:    tt.fields.Day,
			Month:  tt.fields.Month,
			Dow:    tt.fields.Dow,
		}
		nt := time.Now()
		got := s.Next(tt.args.t)
		fmt.Println(got.Sub(nt))
	}
}
