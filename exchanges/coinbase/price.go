package coinbase

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

const (
	tickerKey = "ticker_%s"
)

type (
	Ticker struct {
		TradeID int64           `json:"trade_id"`
		Time    time.Time       `json:"time"`
		Price   decimal.Decimal `json:"price"`
		Bid     decimal.Decimal `json:"bid"`
		Ask     decimal.Decimal `json:"ask"`
		Volume  decimal.Decimal `json:"volume"`
	}
)

func (b *coinbaseEx) getTicker(ctx context.Context, symbol string) (*Ticker, error) {
	cacheKey := fmt.Sprintf(tickerKey, symbol)
	if ticker, ok := b.cache.Get(cacheKey); ok {
		return ticker.(*Ticker), nil
	}

	log := logger.FromContext(ctx)
	uri := fmt.Sprintf("/products/%s/ticker", symbol)
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var ticker Ticker
	if err := UnmarshalResponse(resp, &ticker); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return nil, err
	}

	b.cache.Set(cacheKey, ticker, time.Second*10)
	return &ticker, nil
}
