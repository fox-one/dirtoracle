package bittrex

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

const (
	tickersKey = "tickers"
)

type (
	Ticker struct {
		Symbol        string          `json:"symbol"`
		AskRate       decimal.Decimal `json:"askRate"`
		BidRate       decimal.Decimal `json:"bidRate"`
		LastTradeRate decimal.Decimal `json:"lastTradeRate"`
	}
)

func (b *bittrexEx) getTickers(ctx context.Context) ([]*Ticker, error) {
	if tickers, ok := b.cache.Get(tickersKey); ok {
		return tickers.([]*Ticker), nil
	}

	log := logger.FromContext(ctx)
	uri := "/markets/tickers"
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
