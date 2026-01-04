package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"consistency-lab/internal/metrics"
)

type putBody struct {
	Value string `json:"value"`
}

type putResp struct {
	Item struct {
		ID      string `json:"id"`
		Value   string `json:"value"`
		Version int64  `json:"version"`
	} `json:"item"`
	Mode string `json:"mode"`
}

type item struct {
	ID      string `json:"id"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

func main() {
	var (
		aURL     = flag.String("a", "http://localhost:8081", "node A base url")
		bURL     = flag.String("b", "http://localhost:8082", "node B base url")
		mode     = flag.String("mode", "eventual", "eventual|strong")
		rps      = flag.Int("rps", 50, "requests per second (writes)")
		dur      = flag.Duration("dur", 10*time.Second, "test duration")
		maxLag   = flag.Duration("readlag", 300*time.Millisecond, "max delay between write and read")
		seed     = flag.Int64("seed", 42, "random seed")
		converge = flag.Duration("converge", 2*time.Second, "max time to wait to observe item on both nodes")
	)
	flag.Parse()
	sem := make(chan struct{}, 200)

	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 500,
		IdleConnTimeout:     120 * time.Second,
	}

	client := &http.Client{
		Timeout:   2 * time.Second,
		Transport: tr,
	}

	report := metrics.NewSafeReport(*mode, *dur)

	start := time.Now()
	stopAt := start.Add(*dur)

	interval := time.Second / time.Duration(*rps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	for i := 0; time.Now().Before(stopAt); i++ {
		<-ticker.C
		sem <- struct{}{}
		wg.Add(1)

		workerID := i
		report.IncIterations(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			// each goroutine gets its own RNG (no shared mutable state)
			localRnd := rand.New(rand.NewSource(*seed + int64(workerID)))

			runOne(localRnd, client, *aURL, *bURL, *mode, *maxLag, *converge, report)
		}()
	}
	wg.Wait()

	fmt.Println(report.String())
}

func runOne(rnd *rand.Rand, client *http.Client, aURL, bURL, mode string, maxLag, convergeMax time.Duration, report *metrics.SafeReport) {
	// pick write node
	writeURL := aURL
	if rnd.Intn(2) == 1 {
		writeURL = bURL
	}

	id := fmt.Sprintf("id-%d-%d", time.Now().UnixNano(), rnd.Intn(1_000_000))
	val := fmt.Sprintf("v-%d", rnd.Intn(1_000_000))

	// WRITE
	wStart := time.Now()
	putR, status, err := doPut(client, writeURL, mode, id, val)

	wLat := time.Since(wStart)

	if err != nil {
		report.IncWriteNetErr()
		report.AddWrite(wLat, false)
		return
	}

	report.IncWriteHTTPResp()

	if status/100 == 2 {
		report.IncWriteHit200()
		report.AddWrite(wLat, true)
	} else {
		report.AddWrite(wLat, false)
		return
	}

	// wait random lag then READ from random node
	time.Sleep(time.Duration(rnd.Int63n(int64(maxLag))))

	readURL := aURL
	if rnd.Intn(2) == 1 {
		readURL = bURL
	}

	rStart := time.Now()
	got, status, err := doGet(client, readURL, id)
	rLat := time.Since(rStart)

	if err != nil {
		// Transport error (timeout/EOF/refused). Not a stale read.
		report.IncReadNetErr()
		report.AddRead(rLat, false, false)
	} else {
		report.IncReadHTTPResp()

		switch status {
		case 200:
			report.IncReadHit200()
			stale := got.Version < putR.Item.Version
			report.AddRead(rLat, true, stale)
		case 404:
			// Not found => stale visibility (common in eventual)
			report.AddRead(rLat, false, true)
		default:
			// Other HTTP errors (e.g. 500/503). Not stale.
			report.AddRead(rLat, false, false)
		}
	}

	// convergence: wait until both nodes see >= written version (or timeout)
	if rnd.Float64() > 0.05 { // ~5% sampling
		return
	}

	cStart := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), convergeMax)
	defer cancel()

	for {
		aIt, aStatus, aErr := doGet(client, aURL, id)
		bIt, bStatus, bErr := doGet(client, bURL, id)

		if aErr == nil && bErr == nil &&
			aStatus == 200 && bStatus == 200 &&
			aIt.Value == val && bIt.Value == val {

			report.AddConverge(time.Since(cStart))
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func doPut(client *http.Client, baseURL, mode, id, val string) (putResp, int, error) {
	body := putBody{Value: val}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, baseURL+"/items/"+id, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Consistency", mode)

	resp, err := client.Do(req)
	if err != nil {
		return putResp{}, 0, err // transporte (DNS/EOF/reset/timeout)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return putResp{}, resp.StatusCode, nil
	}

	var pr putResp
	_ = json.NewDecoder(resp.Body).Decode(&pr)
	return pr, resp.StatusCode, nil
}

func doGet(client *http.Client, baseURL, id string) (item, int, error) {
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/items/"+id, nil)
	resp, err := client.Do(req)
	if err != nil {
		return item{}, 0, err // transporte (DNS/EOF/reset/timeout)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return item{}, resp.StatusCode, nil
	}
	var it item
	_ = json.NewDecoder(resp.Body).Decode(&it)
	return it, resp.StatusCode, nil
}
