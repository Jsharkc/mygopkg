package stringutil

import "testing"

func TestSlice(t *testing.T) {
	type args struct {
		str string
		i   int
		j   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "0-5", args: args{str: "这是一个megaview看板", i: 0, j: 5}, want: "这是一个m"},
		{name: "1-5", args: args{str: "这是一个megaview看板", i: 1, j: 5}, want: "是一个m"},
		{name: "0-15", args: args{str: "这是一个megaview看板", i: 0, j: 15}, want: "这是一个megaview看板"},
		{name: "2-15", args: args{str: "这是一个megaview看板", i: 2, j: 15}, want: "一个megaview看板"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slice(tt.args.str, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}
