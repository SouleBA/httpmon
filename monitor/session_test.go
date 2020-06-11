package monitor

import (
	"io/ioutil"
	"testing"
)

func TestCheckAlert(t *testing.T) {
	tests := []struct {
		traffic       []uint
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
			totalEntries: tt.traffic,
			totalPolls:   tt.polls,
		}

		s.checkAlert(ioutil.Discard, 10)

		if s.isAlert != tt.expectedAlert {
			t.Errorf("TestCheckAlert() error = wrong alert: \n\t expected \n%#v \n\t got \n%#v", tt.expectedAlert, s.isAlert)
		}

	}
}

func TestReport(t *testing.T) {
	tests := []struct {
		traffic       []uint
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
			totalEntries: tt.traffic,
			totalPolls:   tt.polls,
		}

		s.checkAlert(ioutil.Discard, 10)

		if s.isAlert != tt.expectedAlert {
			t.Errorf("TestCheckAlert() error = wrong alert: \n\t expected \n%#v \n\t got \n%#v", tt.expectedAlert, s.isAlert)
		}

	}
}
