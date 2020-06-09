package monitor

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestLaunch(t *testing.T) {
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
		setC := func(filepath string, treshold uint, maxPoll uint) options {
			return func(l *Launcher) {
				l.filePath = filepath
				l.treshold = treshold
				l.maxPoll = maxPoll
			}
		}
		l := NewLauncher(setC("fakeDir/access.log", 10, 20))
		afero.WriteFile(appFs, "fakeDir/access.log", []byte(tt.input), 0644)
		l.Launch(os.Stdout)

		fmt.Println(l.session.totalPolls)

		/*for i, val := range e {
			if val != tt.expectedField[i] {
				t.Errorf("parseContent() error: \n\t expected %v \n\t got %v\n", tt.expectedField[i], val)
			}
		}*/

	}

	appFs.RemoveAll("fakeDir")

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
				t.Errorf("parseContent() error: \n\t expected %v \n\t got %v\n", tt.expectedField[i], val)
			}
		}

	}

	appFs.RemoveAll("fakeDir")

}
