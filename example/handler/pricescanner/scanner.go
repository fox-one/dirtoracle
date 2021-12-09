package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/example/config"
	"github.com/fox-one/pkg/uuid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Handle(cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Handle("/", handle(cfg))
	return r
}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	_ = enc.Encode(map[string]interface{}{"data": v})
}

func handle(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requests []*core.PriceRequest
		for _, asset := range cfg.Assets {
			requests = append(requests, &core.PriceRequest{
				TraceID: uuid.Modify(asset.AssetID, fmt.Sprintf("price-request:%s:%d", cfg.Dapp.ClientID, time.Now().Unix()/60)),
				Asset:   *asset,
				Receiver: &core.Receiver{
					Threshold: 1,
					Members:   []string{cfg.Dapp.ClientID},
				},
				Signers:   cfg.Signers,
				Threshold: cfg.Threshold,
			})
		}
		JSON(w, requests)
	}
}
