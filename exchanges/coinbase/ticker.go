package coinbase

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/core"
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

func convertTicker(t *Ticker) *core.Ticker {
	return &core.Ticker{
		Source:    exchangeName,
		Timestamp: t.Time.Unix() * 1000,
		Price:     t.Price,
		VolumeUSD: t.Volume.Mul(t.Price),
	}
}

func readTicker(ctx context.Context, symbol string) (*core.Ticker, error) {
	uri := fmt.Sprintf("/products/%s/ticker", symbol)
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
