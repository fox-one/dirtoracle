package kraken

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Time         time.Time
		Ask          [3]decimal.Decimal `json:"a"`
		Bid          [3]decimal.Decimal `json:"b"`
		Last         [2]decimal.Decimal `json:"c"`
		Volume       [2]decimal.Decimal `json:"v"`
		AveragePrice [2]decimal.Decimal `json:"p"`
		Trades       [2]int             `json:"t"`
		High         [2]decimal.Decimal `json:"h"`
		Low          [2]decimal.Decimal `json:"l"`
		Open         decimal.Decimal    `json:"o"`
	}
)

func convertTicker(t *Ticker) *core.Ticker {
	return &core.Ticker{
		Source:    exchangeName,
		Timestamp: t.Time.Unix() * 1000,
		Price:     t.Last[0],
		VolumeUSD: t.Volume[1].Mul(t.Last[0]),
	}
}

func readTicker(ctx context.Context, symbol string) (*core.Ticker, error) {
	uri := "/public/Ticker"
	resp, err := Request(ctx).SetQueryParam("pair", symbol).Get(uri)
	if err != nil {
		return nil, err
	}

	var ts = map[string]*Ticker{}
	if err := UnmarshalResponse(resp, &ts); err != nil {
		return nil, err
	}

	time, err := time.Parse(time.RFC1123, resp.Header().Get("date"))
	for _, t := range ts {
		t.Time = time
		return convertTicker(t), nil
	}

	return nil, errors.New(string(resp.Body()))
}
