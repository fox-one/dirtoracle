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
		for _, asset := range []core.Asset{
			{AssetID: "c94ac88f-4671-3976-b60a-09064f1811e8", Symbol: "XIN"},
			{AssetID: "c6d0c728-2624-429b-8e0d-d9d19b6592fa", Symbol: "BTC"},
			{AssetID: "6cfe566e-4aad-470b-8c9a-2fd35b49c68d", Symbol: "EOS"},
			{AssetID: "43d61dcd-e413-450d-80b8-101d5e903357", Symbol: "ETH"},
			{AssetID: "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c", Symbol: "BOX"},
		} {
			requests = append(requests, &core.PriceRequest{
				TraceID: uuid.Modify(asset.AssetID, fmt.Sprintf("price-request:%s:%d", cfg.Dapp.ClientID, time.Now().Unix()/600)),
				Asset:   asset,
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
