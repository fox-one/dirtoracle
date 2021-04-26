package bitfinex

import (
	"context"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type (
	Prices map[string]decimal.Decimal
)

func (b *bitfinexEx) getPrices(ctx context.Context) (Prices, error) {
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

	return prices, nil
}

func (b *bitfinexEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	uri := "/tickers"
	resp, err := Request(ctx).SetQueryParam("symbols", symbol).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return decimal.Zero, err
	}

	var tickers [][]interface{}
	if err := UnmarshalResponse(resp, &tickers); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return decimal.Zero, err
	}

	if len(tickers) == 1 && len(tickers[0]) >= 11 {
		if price, ok := tickers[0][7].(float64); ok {
			return decimal.NewFromFloat(price), nil
		}
	}
	return decimal.Zero, nil
}
