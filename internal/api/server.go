package api

import "net/http"

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /items/", s.PutItem)
	mux.HandleFunc("GET /items/", s.GetItem)

	mux.HandleFunc("POST /internal/replicate", s.InternalReplicate)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	return mux
}
