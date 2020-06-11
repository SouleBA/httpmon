package monitor

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestLaunch(t *testing.T) {
	var appFs = afero.NewOsFs()
	appFs.MkdirAll("fakeDir", 0755)
	tests := []struct {
		input                   []string
		runtime                 uint
		expectedSectionHitList  hitList
		expectedProtocolHitList hitList
		expectedStatusHitList   hitList
	}{
		{
			[]string{"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"},
			9,
			[]hit{
				{"/report", 1},
			},
			[]hit{
				{"HTTP/1.0", 1},
			},
			[]hit{
				{"200", 1},
			},
		},
		{
			[]string{"127.0.0.1 - jill [09/May/2018:16:00:41 +0000] \"GET /api/user HTTP/1.0\" 200 234",
				"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"},
			9,
			[]hit{
				{"/api", 1},
				{"/report", 1},
			},
			[]hit{
				{"HTTP/1.0", 2},
			},
			[]hit{
				{"200", 2},
			},
		},
		{
			[]string{"127.0.0.1 - jill [09/May/2018:16:00:41 +0000] \"GET /api/user HTTP/1.0\" 200 234",
				"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"},
			11,
			[]hit{},
			[]hit{},
			[]hit{},
		},
	}

	setC := func(filepath string, maxPoll uint) Options {
		return func(l *Launcher) {
			l.filePath = filepath
			l.maxPoll = maxPoll
		}
	}

	for _, tt := range tests {
		l := NewLauncher(setC("fakeDir/access.log", tt.runtime))
		go func() {
			timer1 := time.NewTimer(4 * time.Second)
			for i, in := range tt.input {
				afero.WriteFile(appFs, "fakeDir/access.log", []byte(in), 0644)

				if i+1 < len(tt.input) {
					<-timer1.C
				}
			}
		}()
		l.Launch(ioutil.Discard)

		sectionHit := rankByHitCount(l.session.pollStats.getSectionList())
		compareSize(t, "section", sectionHit, tt.expectedSectionHitList)
		compareValue(t, "section", sectionHit, tt.expectedSectionHitList)

		protocolHit := rankByHitCount(l.session.pollStats.getProtocolList())
		compareSize(t, "protocol", protocolHit, tt.expectedProtocolHitList)
		compareValue(t, "protocol", protocolHit, tt.expectedProtocolHitList)

		statusHit := rankByHitCount(l.session.pollStats.getStatusList())
		compareSize(t, "status", statusHit, tt.expectedStatusHitList)
		compareValue(t, "status", statusHit, tt.expectedStatusHitList)

	}

	appFs.RemoveAll("fakeDir")

}

func compareValue(t *testing.T, name string, op1 hitList, op2 hitList) {
	for i, val := range op1 {
		if val != op2[i] {
			t.Errorf("TestLaunch() %s hitList error: \n\t expected %v \n\t got %v\n", name, op2[i], val)
		}
	}
}

func compareSize(t *testing.T, name string, op1 hitList, op2 hitList) {
	if len(op1) != len(op2) {
		t.Errorf("TestLaunch() %s hitList error: \n\t expected size %v \n\t got %v\n", name, len(op1), len(op2))
	}
}

func TestParseContent(t *testing.T) {
	var appFs = afero.NewOsFs()
	appFs.MkdirAll("fakeDir", 0755)
	tests := []struct {
		input         string
		expectedField []entry
	}{
		{
			"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123",
			[]entry{{
				remoteIP: "127.0.0.1",
				request: request{
					method:   "GET",
					section:  "/report",
					protocol: "HTTP/1.0",
				},
				responseSize: 123,
				statusCode:   "200",
			}},
		},
		{
			join("127.0.0.1 - jill [09/May/2018:16:00:41 +0000] \"GET /api/user HTTP/1.0\" 200 234",
				"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"),
			[]entry{{
				remoteIP: "127.0.0.1",
				request: request{
					method:   "GET",
					section:  "/api",
					protocol: "HTTP/1.0",
				},
				responseSize: 234,
				statusCode:   "200",
			}, {
				remoteIP: "127.0.0.1",
				request: request{
					method:   "GET",
					section:  "/report",
					protocol: "HTTP/1.0",
				},
				responseSize: 123,
				statusCode:   "200",
			}},
		},
	}

	for _, tt := range tests {
		logReader := strings.NewReader(tt.input)
		l := NewLauncher()
		e, err := l.parseContent(logReader)
		if err != nil {
			t.Errorf("parseContent() error = %v", err)
		}

		for i, val := range e {
			if val != tt.expectedField[i] {
				t.Errorf("TestParseContent() error: \n\t expected %v \n\t got %v\n", tt.expectedField[i], val)
			}
		}

	}

	appFs.RemoveAll("fakeDir")

}
