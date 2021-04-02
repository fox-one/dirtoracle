package bitfinex

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

const (
	pricesKey = "prices"
)

type (
	Prices map[string]decimal.Decimal
)

func (b *bitfinexEx) getPrices(ctx context.Context) (Prices, error) {
	if prices, ok := b.cache.Get(pricesKey); ok {
		return prices.(Prices), nil
	}

	log := logger.FromContext(ctx)
	uri := "/tickers"
	resp, err := Request(ctx).SetQueryParam("symbols", "ALL").Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var tickers [][]interface{}
	if err := UnmarshalResponse(resp, &tickers); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return nil, err
	}

	prices := Prices{}
	for _, d := range tickers {
		if len(d) < 11 {
			continue
		}
		if symbol, ok := d[0].(string); ok {
			if price, ok := d[7].(float64); ok {
				prices[symbol] = decimal.NewFromFloat(price)
			}
		}

	}

	b.cache.Set(pricesKey, prices, time.Second*10)
	return prices, nil
}
