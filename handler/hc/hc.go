package hc

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Handle(version string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Handle("/", handle(version))
	return r
}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	_ = enc.Encode(v)
}

func handle(version string) http.HandlerFunc {
	b := time.Now()
	return func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(b).Truncate(time.Millisecond)
		JSON(w, map[string]interface{}{
			"uptime":  uptime.String(),
			"version": version,
		})
	}
}
