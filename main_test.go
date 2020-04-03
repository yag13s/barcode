package main

import (
	"testing"
)

func Test_getFilename(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			"[success]",
			52,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFilename()
			if len(got) != tt.length {
				t.Errorf("getFilename() = %v, length %v", len(got), tt.length)
			}
		})
	}
}

func Test_randString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"[success] :)",
			args{6},
			6,
		},
		{
			"[success] :)",
			args{-1},
			0,
		},
		{
			"[success] :)",
			args{0},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randString(tt.args.n); len(got) != tt.want {
				t.Errorf("randString() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
