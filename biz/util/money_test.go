package util

import "testing"

func TestParseAmount5DP(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   int64
		hasErr bool
	}{
		{name: "five decimal places", input: "100.23344", want: 10023344},
		{name: "pads fewer decimals", input: "100.23", want: 10023000},
		{name: "whole number", input: "100", want: 10000000},
		{name: "zero", input: "0", want: 0},
		{name: "too many decimals", input: "100.233441", hasErr: true},
		{name: "invalid", input: "abc", hasErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAmount5DP(tt.input)
			if tt.hasErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
