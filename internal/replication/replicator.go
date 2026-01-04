package replication

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"consistency-lab/internal/store"
)

type Replicator struct {
	PeerURL string
	Client  *http.Client
	Chaos   *Chaos
}

type replicateReq struct {
	Item store.Item `json:"item"`
}

func (r *Replicator) ReplicateAsync(_ context.Context, it store.Item) {
	go func() {
		// Detach from request context. Replication is best-effort and should not be
		// canceled just because the client request finished.
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_ = r.Replicate(ctx, it)
	}()
}

// Replicate sends to peer and waits response.
func (r *Replicator) Replicate(ctx context.Context, it store.Item) error {
	if r.PeerURL == "" {
		return errors.New("peer url not set")
	}

	if r.Chaos != nil && r.Chaos.ShouldDrop() {
		// dropped
		return errors.New("replication dropped (chaos)")
	}

	time.Sleep(r.Chaos.Delay())

	payload := replicateReq{Item: it}
	b, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.PeerURL+"/internal/replicate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errors.New("replication failed: non-2xx")
	}
	return nil
}
