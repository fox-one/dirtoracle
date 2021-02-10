package bittrex

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/shopspring/decimal"
)

type (
	Summary struct {
		Symbol        string          `json:"symbol"`
		UpdatedAt     time.Time       `json:"updatedAt"`
		High          decimal.Decimal `json:"high"`
		Low           decimal.Decimal `json:"low"`
		Volume        decimal.Decimal `json:"volume"`
		QuoteVolume   decimal.Decimal `json:"quoteVolume"`
		PercentChange decimal.Decimal `json:"percentChange"`
	}

	Ticker struct {
		Symbol        string          `json:"symbol"`
		AskRate       decimal.Decimal `json:"askRate"`
		BidRate       decimal.Decimal `json:"bidRate"`
		LastTradeRate decimal.Decimal `json:"lastTradeRate"`
	}
)

func convertTicker(s *Summary, t *Ticker) *core.Ticker {
	return &core.Ticker{
		Source:    exchangeName,
		Timestamp: s.UpdatedAt.Unix() * 1000,
		Price:     t.LastTradeRate,
		VolumeUSD: s.QuoteVolume,
	}
}

func readTickers(ctx context.Context) (map[string]*core.Ticker, error) {
	var summaries = map[string]*Summary{}

	{
		uri := "/markets/summaries"
		resp, err := Request(ctx).Get(uri)
		if err != nil {
			return nil, err
		}

		var arr []*Summary
		if err := UnmarshalResponse(resp, &arr); err != nil {
			return nil, err
		}

		for _, s := range arr {
			summaries[s.Symbol] = s
		}
	}

	var tickers = map[string]*core.Ticker{}
	{
		uri := "/markets/tickers"
		resp, err := Request(ctx).Get(uri)
		if err != nil {
			return nil, err
		}

		var arr []*Ticker
		if err := UnmarshalResponse(resp, &arr); err != nil {
			return nil, err
		}

		for _, t := range arr {
			if s, ok := summaries[t.Symbol]; ok {
				tickers[s.Symbol] = convertTicker(s, t)
			}
		}
	}

	return tickers, nil
}
