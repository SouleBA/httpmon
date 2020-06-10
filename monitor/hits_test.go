package monitor

import (
	"testing"
)

func TestRankByHitCount(t *testing.T) {
	tests := []struct {
		sectionHitlist  map[string]int
		expectedHitList []hit
	}{
		{
			map[string]int{
				"/report": 26,
				"/api":    164,
				"/pages":  60,
			},
			[]hit{
				{"/api", 164},
				{"/pages", 60},
				{"/report", 26},
			},
		},
		{
			map[string]int{
				"/report": 43,
				"/api":    4,
				"/pages":  23,
			},
			[]hit{
				{"/report", 43},
				{"/pages", 23},
				{"/api", 4},
			},
		},
	}

	for _, tt := range tests {
		hits := rankByHitCount(tt.sectionHitlist)
		for i, hit := range hits {
			if hit != tt.expectedHitList[i] {
				t.Errorf("TestRankByHitCount() error: \n\t expected \n%#v \n\t got \n%#v", tt.expectedHitList[i], hit)
			}
		}

	}
}
