package main

import "testing"

func TestBuildDateFromEpoch(t *testing.T) {
	originalEpoch, originalDate := BuildEpoch, BuildDate
	t.Cleanup(func() {
		BuildEpoch, BuildDate = originalEpoch, originalDate
	})

	tests := []struct {
		name  string
		epoch string
		want  string
	}{
		{name: "valid epoch", epoch: "1700000000", want: "2023-11-14T22:13:20Z"},
		{name: "unset epoch", epoch: "0", want: "N/A"},
		{name: "negative epoch", epoch: "-1", want: "N/A"},
		{name: "invalid epoch", epoch: "invalid", want: "N/A"},
		{name: "empty epoch", epoch: "", want: "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildEpoch = tt.epoch
			BuildDate = "N/A"
			setBuildDateFromEpoch()
			if BuildDate != tt.want {
				t.Errorf("BuildDate = %q, want %q", BuildDate, tt.want)
			}
		})
	}
}
