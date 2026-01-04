package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"consistency-lab/internal/replication"
	"consistency-lab/internal/store"
)

type Server struct {
	Node       string
	Store      *store.Store
	Replicator *replication.Replicator
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) PutItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		writeJSON(w, 400, ErrResp{Error: "missing id"})
		return
	}

	var body PutBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Value == "" {
		writeJSON(w, 400, ErrResp{Error: "invalid body"})
		return
	}

	mode := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Consistency")))
	if mode == "" {
		mode = "eventual"
	}
	if mode != "eventual" && mode != "strong" {
		writeJSON(w, 400, ErrResp{Error: "X-Consistency must be 'eventual' or 'strong'"})
		return
	}

	it := s.Store.Put(id, body.Value)

	// Context with hard cap to avoid hanging strong writes forever in chaos/partition
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if mode == "eventual" {
		s.Replicator.ReplicateAsync(ctx, it)
		writeJSON(w, 200, PutResp{Item: it, Mode: mode})
		return
	}

	// strong (simulated 2/2): wait for peer ACK
	if err := s.Replicator.Replicate(ctx, it); err != nil {
		writeJSON(w, 503, ErrResp{Error: "strong write failed (peer not ack): " + err.Error()})
		return
	}
	writeJSON(w, 200, PutResp{Item: it, Mode: mode})
}

func (s *Server) GetItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/items/")
	if id == "" {
		writeJSON(w, 400, ErrResp{Error: "missing id"})
		return
	}
	it, ok := s.Store.Get(id)
	if !ok {
		writeJSON(w, 404, ErrResp{Error: "not found"})
		return
	}
	writeJSON(w, 200, it)
}

// internal replication endpoint
func (s *Server) InternalReplicate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Item store.Item `json:"item"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, ErrResp{Error: "invalid body"})
		return
	}
	s.Store.ApplyReplica(req.Item)
	writeJSON(w, 200, map[string]any{"ok": true, "node": s.Node})
}
