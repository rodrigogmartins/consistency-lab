package metrics

import (
	"fmt"
	"time"
)

type Counters struct {
	WriteOK int
	WriteNG int
	ReadOK  int
	ReadNG  int

	// Write outcomes
	WriteHit200   int // status 2xx (you can keep as 2xx)
	WriteHTTPResp int // got an HTTP response (any status)
	WriteNetErr   int // transport error (timeout/EOF/refused)

	// Read outcomes
	ReadHit200   int // status 200
	ReadHTTPResp int // got an HTTP response (any status)
	ReadNetErr   int // transport error

	// Consistency signal (keep your current one)
	StaleRead int
}

type Report struct {
	Mode       string
	WriteLat   Samples
	ReadLat    Samples
	Converge   Samples
	Counters   Counters
	Duration   time.Duration
	Iterations int
}

func (r Report) String() string {
	// ---- WRITE METRICS ----
	writeAttempts := r.Counters.WriteHTTPResp + r.Counters.WriteNetErr

	writeServiceAvail := 0.0
	if writeAttempts > 0 {
		writeServiceAvail = float64(r.Counters.WriteHTTPResp) / float64(writeAttempts) * 100.0
	}

	writeHitRate := 0.0
	if r.Counters.WriteHTTPResp > 0 {
		writeHitRate = float64(r.Counters.WriteHit200) / float64(r.Counters.WriteHTTPResp) * 100.0
	}

	// ---- READ METRICS ----
	readAttempts := r.Counters.ReadHTTPResp + r.Counters.ReadNetErr

	readServiceAvail := 0.0
	if readAttempts > 0 {
		readServiceAvail = float64(r.Counters.ReadHTTPResp) / float64(readAttempts) * 100.0
	}

	readHitRate := 0.0
	if r.Counters.ReadHTTPResp > 0 {
		readHitRate = float64(r.Counters.ReadHit200) / float64(r.Counters.ReadHTTPResp) * 100.0
	}

	staleRate := 0.0
	if readAttempts > 0 {
		staleRate = float64(r.Counters.StaleRead) / float64(readAttempts) * 100.0
	}

	return fmt.Sprintf(
		`MODE: %s
duration: %s | iterations: %d

WRITE
  latency: p50 %s | p95 %s | p99 %s | samples %d
  service_availability: %.2f%% (http_resp=%d net_err=%d)
  write_hit_rate_2xx:   %.2f%% (%d/%d)

READ
  latency: p50 %s | p95 %s | p99 %s | samples %d
  service_availability: %.2f%% (http_resp=%d net_err=%d)
  read_hit_rate_200:    %.2f%% (%d/%d)
  stale_read_rate:      %.2f%% (%d/%d)

CONVERGENCE
  time_to_visibility: p50 %s | p95 %s | samples %d
`,
		r.Mode,
		r.Duration, r.Iterations,

		// WRITE
		r.WriteLat.Percentile(50),
		r.WriteLat.Percentile(95),
		r.WriteLat.Percentile(99),
		r.WriteLat.Count(),
		writeServiceAvail,
		r.Counters.WriteHTTPResp,
		r.Counters.WriteNetErr,
		writeHitRate,
		r.Counters.WriteHit200,
		r.Counters.WriteHTTPResp,

		// READ
		r.ReadLat.Percentile(50),
		r.ReadLat.Percentile(95),
		r.ReadLat.Percentile(99),
		r.ReadLat.Count(),
		readServiceAvail,
		r.Counters.ReadHTTPResp,
		r.Counters.ReadNetErr,
		readHitRate,
		r.Counters.ReadHit200,
		r.Counters.ReadHTTPResp,
		staleRate,
		r.Counters.StaleRead,
		readAttempts,

		// CONVERGENCE
		r.Converge.Percentile(50),
		r.Converge.Percentile(95),
		r.Converge.Count(),
	)
}
