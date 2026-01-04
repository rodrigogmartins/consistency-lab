package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"consistency-lab/internal/api"
	"consistency-lab/internal/replication"
	"consistency-lab/internal/store"
)

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getenvFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func getenvDur(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func main() {
	node := getenv("NODE_NAME", "node_a")
	addr := getenv("ADDR", ":8080")
	peer := getenv("PEER_URL", "")

	seed := time.Now().UnixNano()
	rnd := rand.New(rand.NewSource(seed))

	chaos := replication.Chaos{
		DropRate: getenvFloat("DROP_RATE", 0.0),
		MinDelay: getenvDur("MIN_DELAY", 20*time.Millisecond),
		MaxDelay: getenvDur("MAX_DELAY", 200*time.Millisecond),
		Rand:     rnd,
	}

	s := &api.Server{
		Node:  node,
		Store: store.New(node),
		Replicator: &replication.Replicator{
			PeerURL: peer,
			Client:  &http.Client{Timeout: 2 * time.Second},
			Chaos:   &chaos,
		},
	}

	log.Printf("[%s] listening on %s peer=%s drop=%.2f delay=[%s..%s] seed=%d",
		node, addr, peer, chaos.DropRate, chaos.MinDelay, chaos.MaxDelay, seed)

	if err := http.ListenAndServe(addr, s.Routes()); err != nil {
		log.Fatal(err)
	}
}
