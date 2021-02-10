package bitstamp

import (
	"context"
	"fmt"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/number"
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

func convertTicker(t *Ticker) *core.Ticker {
	ts, _ := number.ParseInt64(t.Timestamp)
	return &core.Ticker{
		Source:    exchangeName,
		Timestamp: ts * 1000,
		Price:     t.Last,
		VolumeUSD: t.Volume.Mul(t.Last),
	}
}

func readTicker(ctx context.Context, symbol string) (*core.Ticker, error) {
	uri := fmt.Sprintf("/ticker/%s", symbol)
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var t Ticker
	if err := UnmarshalResponse(resp, &t); err != nil {
		return nil, err
	}
	return convertTicker(&t), nil
}
