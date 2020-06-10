package monitor

import (
	"fmt"
	"io"
	"time"
)

type session struct {
	pollStats
	traffic
	offset       int64
	fileSize     int64
	isAlert      bool
	recoveryTime time.Time
}

func newSession() *session {
	s := &session{
		pollStats{
			sectionList:  make(map[string]int),
			statusList:   make(map[string]int),
			protocolList: make(map[string]int),
		},
		traffic{},
		0,
		0,
		false,
		time.Now(),
	}

	return s
}

func (s *session) resetPollStats() {
	s.pollStats = pollStats{
		sectionList:  make(map[string]int),
		statusList:   make(map[string]int),
		protocolList: make(map[string]int),
	}
}

func (s *session) resetTraffic() {
	s.traffic = traffic{}
}

func (s *session) getAlertStatus() bool {
	return s.isAlert
}

func (s *session) updateTotalTraffic(entries uint) {
	s.traffic.appendEntries(entries)
}

func (s *session) updateTotalPolls(polls uint) {
	s.traffic.incrementPolls(polls)
}

func (s *session) updateSectionStats(key string) {
	s.pollStats.incrementSectionList(key)
}

func (s *session) updateStatusStats(key string) {
	s.pollStats.incrementStatusList(key)
}

func (s *session) updateProtocolStats(key string) {
	s.pollStats.incrementProtocolList(key)
}

func (s *session) updateOffset(offset int64) {
	s.offset = offset
}

func (s *session) updateFileSize(fileSize int64) {
	s.fileSize = fileSize
}

func (s *session) getOffset() int64 {
	return s.offset
}

func (s *session) getFileSize() int64 {
	return s.fileSize
}

func (s *session) updateRecoveryDate(recovery time.Time) {
	s.recoveryTime = recovery
}

func (s *session) getTrafficAvg() uint {
	return uint(s.traffic.hitsAvgRate())
}

func (s *session) checkAlert(out io.Writer, treshold uint) {
	tAvg := uint(s.traffic.hitsAvgRate())
	if !s.isAlert && tAvg > treshold {
		s.isAlert = true
		s.notify(out, tAvg)
	}

	if s.isAlert && tAvg < treshold {
		s.isAlert = false
		s.notify(out, tAvg)
	}
}

func (s *session) notify(out io.Writer, tAvg uint) {
	if s.isAlert {
		fmt.Fprintf(out, "High traffic generated an alert - hits = %d\n", tAvg)
		fmt.Fprintf(out, "\n")
		return
	}

	if !s.isAlert {
		fmt.Fprintf(out, "High traffic recovered at %s - hits = %d\n", s.recoveryTime, tAvg)
		fmt.Fprintf(out, "\n")
		return
	}
}

func (s *session) report(out io.Writer, n int) {
	fmt.Fprintln(out, time.Now())
	fmt.Fprintln(out, "Here is your Interval Stats Report", time.Now())
	fmt.Fprintf(out, "\n")

	secStats := rankByHitCount(s.pollStats.getSectionList())
	print("section", secStats, n, out)

	protoStats := rankByHitCount(s.pollStats.getProtocolList())
	print("protocol", protoStats, n, out)

	statusStats := rankByHitCount(s.pollStats.getStatusList())
	print("status", statusStats, n, out)
}

func print(name string, hits hitList, n int, out io.Writer) {
	if len(hits) == 0 {
		fmt.Fprintf(out, "No %s hit for this poll interval\n\n", name)

		return
	} else if n > len(hits) {
		n = len(hits)
	}
	fmt.Fprintf(out, "%s\thits\n", name)
	for i := 0; i < n; i++ {
		fmt.Fprintf(out, "%s\t%d\t\n", hits[i].key, hits[i].value)
	}
	fmt.Fprint(out, "\n")
}
