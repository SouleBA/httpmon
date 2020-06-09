package monitor

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/satyrius/gonx"
)

// Launcher holds the needed information in order to monitoring
type Launcher struct {
	filePath            string   // filePath is the name of the file to monitor
	treshold            uint     // treshold represents the limit to trigger alerts
	logFormat           string   // logFormat represents the log formatting
	logKeys             []string //logKeys represents the formatted fields keys
	pollInterval        uint     // pollInterval represents the file content polling interval
	statsReportInterval uint
	alertingInterval    uint
	isAlert             bool      // specify if we are currently on alert
	stopCh              chan bool // stopCh is the channel which when closed, will shutdown the launcher
}

type options func(*Launcher)

// NewLauncher returns a configured Launcher
// opts is variadic and represents funtional options
func NewLauncher(opts ...options) *Launcher {
	l := &Launcher{
		filePath:            "/tmp/access.log",
		treshold:            10,
		logFormat:           "$remote_addr $user_identifier $userid [$time_local] \"$request\" $status_code $response_size",
		logKeys:             []string{"remote_addr", "response_size", "status_code", "request"},
		pollInterval:        1,
		statsReportInterval: 10,
		alertingInterval:    120,
		isAlert:             false,
		stopCh:              make(chan bool, 1),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// SetConfig is a functional option to configure the Launcher
func SetConfig(filepath string, treshold uint) options {
	return func(l *Launcher) {
		l.filePath = filepath
		l.treshold = treshold
	}
}

// Launch willl start the monitoring service
// which is actually a loop that polls logs at specified time interval
func (l *Launcher) Launch(out io.Writer) {
	pollTicker := time.NewTicker(time.Duration(l.pollInterval) * time.Second)
	statsTicker := time.NewTicker(time.Duration(l.statsReportInterval) * time.Second)

	session := newSession()

	polls := uint(0) // counter to check if it is time for an alert

	for {
		select {
		case <-pollTicker.C:
			polls++ // increment the counter
			content := content{filePath: l.filePath}
			err := content.sync(session.getOffset(), session.getFileSize())
			if err != nil {
				panic(err)
			}

			// set offset and filesize for next poll
			session.updateOffset(content.offset)
			session.updateFileSize(content.fileSize)

			entries, err := l.parseContent(strings.NewReader(content.fields))
			if err != nil {
				panic(err)
			}

			for _, entry := range entries {
				session.updateSectionStats(entry.request.section)
				session.updateStatusStats(entry.statusCode)
				session.updateProtocolStats(entry.request.protocol)
			}

			if len(entries) > 0 {
				session.updateTotalTraffic(uint(len(entries)))
			}

			session.updateTotalPolls(l.pollInterval)

			if polls%l.alertingInterval == 0 {
				session.checkAlert(l.treshold)
			}
		case <-statsTicker.C:
			session.report()
			session.resetPollStats()
		case <-l.stopCh:
			fmt.Fprintf(out, "Shut down requested")
			pollTicker.Stop()
			statsTicker.Stop()
			return
		default:

		}
		/*
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program)
			if evaluated != nil {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}*/
	}
}

func (l *Launcher) parseContent(logReader io.Reader) ([]entry, error) {
	var e []entry
	fields := make(map[string]string, len(l.logKeys))
	reader := gonx.NewReader(logReader, l.logFormat)
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		for _, key := range l.logKeys {
			field, err := rec.Field(key)
			if err != nil {
				return nil, err
			}
			fields[key] = field
		}

		req, err := NewRequest(fields["request"])
		if err != nil {
			return nil, err
		}

		responseSize, err := strconv.ParseUint(fields["response_size"], 10, 64)
		if err != nil {
			return nil, err
		}

		statusCode := fields["status_code"]

		e = append(e, entry{
			remoteIP:     fields["remote_addr"],
			request:      *req,
			responseSize: responseSize,
			statusCode:   statusCode,
		})

	}

	return e, nil

}

// Shutdown send a signal to stop ongoing monitoring tasks
func (l *Launcher) Shutdown() {
	l.stopCh <- true
}

func printErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some problem!\n")
	io.WriteString(out, "launch errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
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
