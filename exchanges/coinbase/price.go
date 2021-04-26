package coinbase

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
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

func (b *coinbaseEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	uri := fmt.Sprintf("/products/%s/ticker", symbol)
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

	return ticker.Bid, nil
}
