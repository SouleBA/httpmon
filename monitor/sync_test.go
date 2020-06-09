package monitor

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestSync(t *testing.T) {
	var appFs = afero.NewOsFs()
	appFs.MkdirAll("fakeDir", 0755)
	tests := []struct {
		input          string
		expectedField  string
		expectedOffset int64
	}{
		{
			"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123",
			"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123",
			77,
		},
		{
			join("127.0.0.1 - jill [09/May/2018:16:00:41 +0000] \"GET /api/user HTTP/1.0\" 200 234",
				"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"),
			join("127.0.0.1 - jill [09/May/2018:16:00:41 +0000] \"GET /api/user HTTP/1.0\" 200 234",
				"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 123"),
			156,
		},
	}

	for _, tt := range tests {
		c := content{
			filePath: "fakeDir/testSync",
		}

		afero.WriteFile(appFs, "fakeDir/testSync", []byte(tt.input), 0644)

		err := c.sync(0, 0)

		if err != nil {
			t.Errorf("sync() error = %v", err)
		}

		if c.fields != tt.expectedField {
			t.Errorf("sync() error = wrong Lines: \n\t expected \n%#v \n\t got \n%#v", tt.expectedField, c.fields)
		}

		if c.offset != tt.expectedOffset {
			t.Errorf("sync() error = wrong offset: \n\t expected %d \n\t got %d", tt.expectedOffset, c.offset)
		}

	}

	appFs.RemoveAll("fakeDir")
}
func TestRead(t *testing.T) {
	c := content{}
	if err := c.read(strings.NewReader("This is test string 1.\nThis is test string 2.")); (err != nil) != false {
		t.Errorf("read() error = %v", err)
	}
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
		sb.WriteString("\n")
	}
	return strings.Trim(sb.String(), "\n")
}
