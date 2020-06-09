package monitor

import "fmt"

type session struct {
	pollStats
	traffic
	offset   int64
	fileSize int64
	isAlert  bool
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

func (s *session) checkAlert(treshold uint) {
	tAvg := uint(s.traffic.hitsAvgRate())
	if !s.isAlert && tAvg > treshold {
		s.isAlert = true
		s.notify(tAvg)
	}

	if s.isAlert && tAvg < treshold {
		s.isAlert = false
		s.notify(tAvg)
	}
}

func (s *session) notify(tAvg uint) {
	if s.isAlert {
		fmt.Println("High traffic generated an alert - hits =", tAvg)
		return
	}

	if !s.isAlert {
		fmt.Println("High traffic recovered - hits =", tAvg)
		return
	}
}

func (s *session) report() {
	fmt.Println(s.traffic.sum())
	fmt.Println(s.traffic.totalPolls)
	fmt.Println(s.traffic.hitsAvgRate())
	fmt.Println(rankByHitCount(s.pollStats.getSectionList()))
	fmt.Println(rankByHitCount(s.pollStats.getProtocolList()))
	fmt.Println(rankByHitCount(s.pollStats.getStatusList()))
}
