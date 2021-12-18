package scanner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fox-one/dirtoracle/apps/mvm/core"
	"github.com/fox-one/pkg/uuid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Handle(
	system core.System,
	assets core.AssetStore,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Handle("/", handle(system, assets))
	return r
}

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	_ = enc.Encode(map[string]interface{}{"data": v})
}

func Error(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	respBody := "internal server error"

	w.Header().Set("Content-Type", "application/json") // Error responses are always JSON
	w.Header().Set("Content-Length", strconv.Itoa(len(respBody)))
	w.WriteHeader(statusCode) // set HTTP status code and send response

	w.Write([]byte(respBody))
}

func handle(
	system core.System,
	assets core.AssetStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assets, err := assets.List(r.Context())
		if err != nil {
			Error(w, err)
			return
		}

		var requests = make([]*core.PriceRequest, 0, len(assets))
		for _, asset := range assets {
			if asset.PriceUpdatedAt != nil && time.Since(*asset.PriceUpdatedAt) < time.Duration(asset.PriceDuration)*time.Second {
				continue
			}

			trace := uuid.Modify(asset.AssetID, fmt.Sprintf("price-request:%s:%d", system.ClientID, time.Now().Unix()/asset.PriceDuration))
			requests = append(requests, &core.PriceRequest{
				TraceID: trace,
				Asset: core.Asset{
					AssetID: asset.AssetID,
					Symbol:  asset.Symbol,
				},
				Receiver: &core.Receiver{
					Threshold: 1,
					Members:   []string{system.ClientID},
				},
				Signers:   system.Signers,
				Threshold: system.SignerThreshold,
			})
		}
		JSON(w, requests)
	}
}
