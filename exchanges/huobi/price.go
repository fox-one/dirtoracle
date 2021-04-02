package huobi

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
)

const (
	tickersKey = "tickers"
)

type (
	Ticker struct {
		Symbol string  `json:"symbol"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
		Volume float64 `json:"vol"`
	}
)

func (b *huobiEx) getTickers(ctx context.Context) ([]*Ticker, error) {
	if tickers, ok := b.cache.Get(tickersKey); ok {
		return tickers.([]*Ticker), nil
	}

	log := logger.FromContext(ctx)
	uri := "/market/tickers"
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var tickers []*Ticker
	if err := UnmarshalResponse(resp, &tickers); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return nil, err
	}

	b.cache.Set(tickersKey, tickers, time.Second*10)
	return tickers, nil
}
