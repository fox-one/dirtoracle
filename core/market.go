package core

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		AssetID   string          `json:"asset_id"`
		Source    string          `json:"source,omitempty"`
		Timestamp int64           `json:"timestamp,omitempty"`
		Price     decimal.Decimal `json:"price,omitempty"`
		VolumeUSD decimal.Decimal `json:"volume_usd,omitempty"`
	}

	MarketStore interface {
		// ticker
		SaveTicker(ctx context.Context, ticker *Ticker) error
		FindTickers(ctx context.Context, assetID string) ([]*Ticker, error)
		AggregateTickers(ctx context.Context, assetID string) (*Ticker, error)
	}
)

func (t *Ticker) ExportProposal() *PriceProposal {
	return &PriceProposal{
		PriceData: PriceData{
			AssetID:   t.AssetID,
			Timestamp: t.Timestamp,
			Price:     t.Price,
		},
	}
}
