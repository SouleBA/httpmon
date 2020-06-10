package monitor

import (
	"os"
	"testing"
)

func TestCheckAlert(t *testing.T) {
	tests := []struct {
		Traffic       []uint
		polls         uint
		expectedAlert bool
	}{
		{
			[]uint{8, 16, 32, 64, 128},
			5,
			true,
		},
		{
			[]uint{8, 8, 8, 8},
			4,
			false,
		},
	}

	s := newSession()
	for _, tt := range tests {
		s.traffic = traffic{
			totalEntries: tt.Traffic,
			totalPolls:   tt.polls,
		}

		s.checkAlert(os.Stdout, 10)

		if s.isAlert != tt.expectedAlert {
			t.Errorf("TestCheckAlert() error = wrong alert: \n\t expected \n%#v \n\t got \n%#v", tt.expectedAlert, s.isAlert)
		}

	}
}
