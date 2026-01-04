package metrics

import (
	"sync"
	"time"
)

// SafeReport wraps Report with a mutex so you can update it from many goroutines safely.
type SafeReport struct {
	mu sync.Mutex
	r  Report
}

func NewSafeReport(mode string, duration time.Duration) *SafeReport {
	return &SafeReport{
		r: Report{
			Mode:     mode,
			Duration: duration,
		},
	}
}

func (s *SafeReport) IncIterations(n int) {
	s.mu.Lock()
	s.r.Iterations += n
	s.mu.Unlock()
}

// ---- Latency + outcome helpers ----

func (s *SafeReport) AddWrite(lat time.Duration, ok bool) {
	s.mu.Lock()
	s.r.WriteLat.Add(lat)
	if ok {
		s.r.Counters.WriteOK++
	} else {
		s.r.Counters.WriteNG++
	}
	s.mu.Unlock()
}

func (s *SafeReport) AddRead(lat time.Duration, ok bool, stale bool) {
	s.mu.Lock()
	s.r.ReadLat.Add(lat)
	if ok {
		s.r.Counters.ReadOK++
	} else {
		s.r.Counters.ReadNG++
	}
	if stale {
		s.r.Counters.StaleRead++
	}
	s.mu.Unlock()
}

func (s *SafeReport) AddConverge(lat time.Duration) {
	s.mu.Lock()
	s.r.Converge.Add(lat)
	s.mu.Unlock()
}

// ---- Service availability / hit rate counters ----
// These should be incremented in the client code based on (status, err).

func (s *SafeReport) IncWriteHTTPResp() {
	s.mu.Lock()
	s.r.Counters.WriteHTTPResp++
	s.mu.Unlock()
}

func (s *SafeReport) IncWriteNetErr() {
	s.mu.Lock()
	s.r.Counters.WriteNetErr++
	s.mu.Unlock()
}

func (s *SafeReport) IncWriteHit200() {
	s.mu.Lock()
	s.r.Counters.WriteHit200++
	s.mu.Unlock()
}

func (s *SafeReport) IncReadHTTPResp() {
	s.mu.Lock()
	s.r.Counters.ReadHTTPResp++
	s.mu.Unlock()
}

func (s *SafeReport) IncReadNetErr() {
	s.mu.Lock()
	s.r.Counters.ReadNetErr++
	s.mu.Unlock()
}

func (s *SafeReport) IncReadHit200() {
	s.mu.Lock()
	s.r.Counters.ReadHit200++
	s.mu.Unlock()
}

// ---- Output ----

func (s *SafeReport) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.r.String()
}
