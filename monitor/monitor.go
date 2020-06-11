package monitor

import (
	"fmt"
	"io"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/satyrius/gonx"
)

// Launcher holds the needed information in order to monitoring
type Launcher struct {
	filePath            string    // filePath is the name of the file to monitor
	treshold            uint      // treshold represents the limit to trigger alerts
	logFormat           string    // logFormat represents the log formatting
	logKeys             []string  // logKeys represents the formatted fields keys
	pollInterval        uint      // pollInterval represents the file content polling interval
	statsReportInterval uint      // Traffic stats Reporting Interval
	alertingInterval    uint      // Traffic hits rate reporting interval
	maxPoll             uint      // maximum possible polls
	session             *session  // session holds all data
	isAlert             bool      // specify if we are currently on alert
	stopCh              chan bool // stopCh is the channel which when closed, will shutdown the launcher
}

// Options is a Launcher configuration function
type Options func(*Launcher)

// NewLauncher returns a configured Launcher
// opts is variadic and represents funtional options
func NewLauncher(opts ...Options) *Launcher {
	l := &Launcher{
		filePath:            "/tmp/access.log",
		treshold:            10,
		logFormat:           "$remote_addr $user_identifier $userid [$time_local] \"$request\" $status_code $response_size",
		logKeys:             []string{"remote_addr", "response_size", "status_code", "request"},
		pollInterval:        1,
		statsReportInterval: 10,
		alertingInterval:    120,
		maxPoll:             1<<bits.UintSize - 1,
		session:             newSession(),
		isAlert:             false,
		stopCh:              make(chan bool, 1),
	}

	// Apply configuration
	for _, opt := range opts {
		opt(l)
	}

	// check consistency
	if l.statsReportInterval > l.pollInterval {
		fmt.Println("reporting interval must be higher than polling interval")
		os.Exit(1)
	}

	if l.maxPoll > l.pollInterval {
		fmt.Println("maximum polls must be higher than polling interval")
		os.Exit(1)
	}

	return l
}

// DefaultConfig is a functional option to configure the Launcher
func DefaultConfig(filepath string, treshold uint) Options {
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

	//session := newSession()

	polls := uint(0) // counter to check if it is time for an alert

	for {
		select {
		case <-pollTicker.C:
			polls++ // increment the counter
			content := content{filePath: l.filePath}
			err := content.sync(l.session.getOffset(), l.session.getFileSize())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// set offset and filesize for next poll
			l.session.updateOffset(content.offset)
			l.session.updateFileSize(content.fileSize)

			// parse polled content
			entries, err := l.parseContent(strings.NewReader(content.fields))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// update stats with current polled content
			for _, entry := range entries {
				l.session.updateSectionStats(entry.request.section)
				l.session.updateStatusStats(entry.statusCode)
				l.session.updateProtocolStats(entry.request.protocol)
			}

			// only update if there was traffic
			if len(entries) > 0 {
				l.session.updateTotalTraffic(uint(len(entries)))
			}

			// If we are on alert and under treshold persit recovery time
			if l.session.getAlertStatus() && l.session.getTrafficAvg() < l.treshold {
				l.session.updateRecoveryDate(time.Now())
			}

			l.session.updateTotalPolls(l.pollInterval)

			// Every alerting Interval
			if polls%l.alertingInterval == 0 {
				l.session.checkAlert(out, l.treshold)
				l.session.resetTraffic()
			}

			// If maximum polls reached we quit
			if polls == l.maxPoll {
				l.stopCh <- true
			}
		case <-statsTicker.C:
			// Report Interval
			l.session.report(out, 5)
			// Reset Stats Data
			l.session.resetPollStats()
		case <-l.stopCh:
			fmt.Fprintf(out, "Shut down requested\nBye!\n")
			pollTicker.Stop()
			statsTicker.Stop()
			return
		}
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
			return nil, errors.Wrap(err, "gonx read failed")
		}

		for _, key := range l.logKeys {
			field, err := rec.Field(key)
			if err != nil {
				return nil, errors.Wrap(err, "gonx could not retrieve a field")
			}
			fields[key] = field
		}

		req, err := newRequest(fields["request"])
		if err != nil {
			return nil, errors.Wrap(err, "could not create a request")
		}

		responseSize, err := strconv.ParseUint(fields["response_size"], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "Type casting string to uint failed")
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
