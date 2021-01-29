package binance

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Timestamp   int64           `json:"E"`
		Event       string          `json:"e"`
		Symbol      string          `json:"s"`
		Close       decimal.Decimal `json:"c"`
		Open        decimal.Decimal `json:"o"`
		High        decimal.Decimal `json:"h"`
		Low         decimal.Decimal `json:"l"`
		BaseVolume  decimal.Decimal `json:"v"`
		QuoteVolume decimal.Decimal `json:"q"`
	}
)

func convertTicker(assetID string, t *Ticker) *core.Ticker {
	return &core.Ticker{
		AssetID:   assetID,
		Source:    exchangeName,
		Timestamp: t.Timestamp,
		Price:     t.Close,
		VolumeUSD: t.QuoteVolume,
	}
}
