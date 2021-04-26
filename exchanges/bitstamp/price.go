package bitstamp

import (
	"context"
	"fmt"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Timestamp string          `json:"timestamp"`
		High      decimal.Decimal `json:"high"`
		Low       decimal.Decimal `json:"low"`
		Last      decimal.Decimal `json:"last"`
		Open      decimal.Decimal `json:"open"`
		VWAP      decimal.Decimal `json:"vwap"`
		Bid       decimal.Decimal `json:"bid"`
		Ask       decimal.Decimal `json:"ask"`
		Volume    decimal.Decimal `json:"volume"`
	}
)

func (b *bitstampEx) getTicker(ctx context.Context, symbol string) (*Ticker, error) {
	log := logger.FromContext(ctx)
	uri := fmt.Sprintf("/ticker/%s", symbol)
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

	return &ticker, nil
}
