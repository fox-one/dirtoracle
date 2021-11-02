package bittrex

import (
	"context"
	"fmt"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Symbol        string          `json:"symbol"`
		AskRate       decimal.Decimal `json:"askRate"`
		BidRate       decimal.Decimal `json:"bidRate"`
		LastTradeRate decimal.Decimal `json:"lastTradeRate"`
	}
)

func (exch *bittrexEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	uri := fmt.Sprintf("/markets/%s/ticker", symbol)
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return decimal.Zero, err
	}

	var ticker Ticker
	if err := UnmarshalResponse(resp, &ticker); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return decimal.Zero, err
	}

	return ticker.BidRate, nil
}
