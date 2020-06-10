package monitor

import "sort"

type hit struct {
	key   string
	value int
}

type hitList []hit

func rankByHitCount(hitFrequencies map[string]int) hitList {
	hl := make(hitList, len(hitFrequencies))
	i := 0
	for k, v := range hitFrequencies {
		hl[i] = hit{k, v}
		i++
	}
	sort.Sort(sort.Reverse(hl))
	return hl
}

func (h hitList) Len() int           { return len(h) }
func (h hitList) Less(i, j int) bool { return h[i].value < h[j].value }
func (h hitList) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
