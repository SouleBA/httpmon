package monitor

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// content represents data retrieved when polling a file
type content struct {
	fields   string // fields are the file's polled fields
	fileSize int64  // fileSize is the polled file size
	filePath string // filePath is the polled file path
	offset   int64  // offset is the offset at the end of the polling
}

func (c *content) sync(initialOffset int64, initialFileSize int64) error {
	stat, err := os.Stat(c.filePath)
	if err != nil {
		return errors.Wrap(err, "Getting FileInfo failed")
	}

	c.fileSize = stat.Size()
	if c.fileSize != initialFileSize {
		f, err := os.Open(c.filePath)
		if err != nil {
			return errors.Wrap(err, "open failed")
		}
		defer f.Close()

		_, err = f.Seek(initialOffset, 0)
		if err != nil {
			return errors.Wrap(err, "set file offset failed")
		}

		err = c.read(f)
		if err != nil {
			return errors.Wrap(err, "read file failed")
		}

		c.offset, err = f.Seek(0, 1)
		if err != nil {
			return errors.Wrap(err, "set file offset failed")
		}
	}

	return nil
}

func (c *content) read(handle io.Reader) error {
	var sb strings.Builder
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		c.append(&sb, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "Scanner encountered a non-EOF error")
	}

	c.fields = strings.Trim(sb.String(), "\n")

	return nil
}

func (c *content) append(sb *strings.Builder, str string) {
	sb.WriteString(str)
	sb.WriteString("\n")

}
