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
				hit{"/api", 164},
				hit{"/pages", 60},
				hit{"/report", 26},
			},
		},
		{
			map[string]int{
				"/report": 43,
				"/api":    4,
				"/pages":  23,
			},
			[]hit{
				hit{"/report", 43},
				hit{"/pages", 23},
				hit{"/api", 4},
			},
		},
	}

	for _, tt := range tests {
		hits := rankByHitCount(tt.sectionHitlist)
		for i, hit := range hits {
			if hit != tt.expectedHitList[i] {
				t.Errorf("rankByHitCount() error = wrong order: \n\t expected \n%#v \n\t got \n%#v", tt.expectedHitList[i], hit)
			}
		}

	}
}
