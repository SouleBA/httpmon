package monitor

import (
	"testing"
)

func TestHitsAvgRate(t *testing.T) {
	tests := []struct {
		Traffic             []uint
		polls               uint
		expectedHitsAvgRate float64
	}{
		{
			[]uint{8, 8, 8, 8},
			4,
			8,
		},
		{
			[]uint{8, 16, 32, 64, 128},
			5,
			50,
		},
	}

	for _, tt := range tests {
		tr := traffic{
			totalEntries: tt.Traffic,
			totalPolls:   tt.polls,
		}

		rate := tr.hitsAvgRate()
		if rate != tt.expectedHitsAvgRate {
			t.Errorf("TestHitsAvgRate error = wrong rate: \n\t expected \n%#v \n\t got \n%#v", tt.expectedHitsAvgRate, rate)
		}

	}
}
