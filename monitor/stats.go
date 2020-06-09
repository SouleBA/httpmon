package monitor

type rank map[string]int

type pollStats struct {
	sectionList  rank
	statusList   rank
	protocolList rank
}

func (p *pollStats) incrementSectionList(key string) {
	increment(p.sectionList, key)
}

func (p *pollStats) getSectionList() rank {
	return p.sectionList
}

func (p *pollStats) incrementProtocolList(key string) {
	increment(p.protocolList, key)

}

func (p *pollStats) getProtocolList() rank {
	return p.protocolList
}

func (p *pollStats) incrementStatusList(key string) {
	increment(p.statusList, key)
}

func increment(r rank, key string) {
	r[key] += 1
}

func (p *pollStats) getStatusList() rank {
	return p.statusList
}

type traffic struct {
	totalEntries []uint
	totalPolls   uint
}

func (t *traffic) sum() uint {
	sum := uint(0)
	for _, val := range t.totalEntries {
		sum += val
	}
	return sum
}

func (t *traffic) hitsAvgRate() float64 {
	return float64(t.sum()) / float64(t.totalPolls)
}

func (t *traffic) appendEntries(entries uint) {
	t.totalEntries = append(t.totalEntries, entries)
}

func (t *traffic) incrementPolls(polls uint) {
	t.totalPolls += polls
}
