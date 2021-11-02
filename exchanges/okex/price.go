package okex

import (
	"context"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Type   string          `json:"instType"`
		Symbol string          `json:"instId"`
		Last   decimal.Decimal `json:"last"`
		Ask    decimal.Decimal `json:"askPx"`
		Bid    decimal.Decimal `json:"bidPx"`
	}
)

func (exch *okexEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	uri := "/v5/market/ticker"
	resp, err := Request(ctx).SetQueryParam("instId", symbol).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return decimal.Zero, err
	}

	var tickers []*Ticker
	if err := UnmarshalResponse(resp, &tickers); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return decimal.Zero, err
	}

	if len(tickers) > 0 {
		return tickers[0].Bid, nil
	}

	return decimal.Zero, nil
}
