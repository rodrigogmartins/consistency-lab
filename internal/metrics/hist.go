package metrics

import (
	"sort"
	"time"
)

type Samples struct {
	values []time.Duration
}

func (s *Samples) Add(d time.Duration) {
	s.values = append(s.values, d)
}

func (s *Samples) Count() int { return len(s.values) }

func (s *Samples) Percentile(p float64) time.Duration {
	if len(s.values) == 0 {
		return 0
	}
	cp := append([]time.Duration(nil), s.values...)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })

	if p <= 0 {
		return cp[0]
	}
	if p >= 100 {
		return cp[len(cp)-1]
	}
	// nearest-rank
	rank := int((p/100.0)*float64(len(cp)-1) + 0.5)
	if rank < 0 {
		rank = 0
	}
	if rank >= len(cp) {
		rank = len(cp) - 1
	}
	return cp[rank]
}
