package main

import "testing"

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
