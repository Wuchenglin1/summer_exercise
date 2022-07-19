package clock

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetField(t *testing.T) {
	type args struct {
		field string
		b     bounds
	}
	tests := []struct {
		name    string
		args    args
		wantS   []Schedule
		wantErr bool
	}{
		{"test minute", args{field: "*/5", b: _minute}, []Schedule{{start: 0, end: 59, hasStep: true, step: 5}}, false},
		{"test minute", args{field: "*-5/5", b: _minute}, []Schedule{{start: 0, end: 59, hasStep: true, step: 5}}, false},
		{"test minute", args{field: "10/8", b: _minute}, []Schedule{{start: 10, end: 59, hasStep: true, step: 8}}, false},
		{"test minute", args{field: "5-20/4", b: _minute}, []Schedule{{start: 5, end: 20, hasStep: true, step: 4}}, false},
	}
	for _, tt := range tests {
		gotS, err := GetField(tt.args.field, tt.args.b)
		fmt.Println(gotS, err)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. GetField() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(gotS, tt.wantS) {
			t.Errorf("%q. GetField() = %v, want %v", tt.name, gotS, tt.wantS)
		}
	}
}

func TestParser_ParseSep(t *testing.T) {
	type fields struct {
		IsWithSecond bool
	}
	type args struct {
		sep string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SepSchedule
		wantErr bool
	}{
		{"test1", fields{IsWithSecond: false}, args{sep: "*/5 1-5/4 1-5/4,4-7/3 * *"}, SepSchedule{
			//Second: []Schedule{{start: 0,end: 59,hasStep:false ,step:0 }},
			Minute: []Schedule{{0, 59, true, 5}},
			Hour:   []Schedule{{1, 5, true, 4}},
			Day:    []Schedule{{1, 5, true, 4}, {4, 7, true, 3}},
			Month:  []Schedule{{1, 12, false, 1}},
			Dow:    []Schedule{{1, 7, false, 1}},
		}, false},
		{"test2", fields{IsWithSecond: false}, args{sep: "1-60/5 1-5/4 1-5/4,4-7/3 * *"}, SepSchedule{
			//Second: []Schedule{{start: 0,end: 59,hasStep:false ,step:0 }},
			Minute: nil,
			Hour:   []Schedule{{1, 5, true, 4}},
			Day:    []Schedule{{1, 5, true, 4}, {4, 7, true, 3}},
			Month:  []Schedule{{1, 12, false, 1}},
			Dow:    []Schedule{{1, 7, false, 1}},
		}, false},
		{"test1", fields{IsWithSecond: false}, args{sep: "6-5/5 0/4 0-5/4,7-70/3 -- *"}, SepSchedule{
			//Second: []Schedule{{start: 0,end: 59,hasStep:false ,step:0 }},
			Minute: nil,
			Hour:   nil,
			Day:    nil,
			Month:  nil,
			Dow:    nil,
		}, true},
	}
	for _, tt := range tests {
		p := Parser{
			IsWithSecond: tt.fields.IsWithSecond,
		}
		got, err := p.ParseSep(tt.args.sep)
		fmt.Println(got, err)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Parser.ParseSep() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Parser.ParseSep() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_parseInt(t *testing.T) {
	type args struct {
		exp string
	}
	tests := []struct {
		name    string
		args    args
		want    uint
		wantErr bool
	}{
		{"test1", args{"10"}, 10, false},
		{"test2", args{"100"}, 100, false},
		{"test3", args{"-1"}, 0, true},
		{"test4", args{"a"}, 0, true},
	}
	for _, tt := range tests {
		got, err := parseInt(tt.args.exp)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. parseInt() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. parseInt() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
