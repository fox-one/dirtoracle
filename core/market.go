package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Exchange string `json:"exchange,omitempty"`
		AssetID  string `json:"asset_id,omitempty"`

		UpdatedAt time.Time       `json:"updated_at,omitempty"`
		Price     decimal.Decimal `json:"last_price,omitempty"`
		VolumeUSD decimal.Decimal `json:"volume_usd,omitempty"`
	}

	MarketStore interface {
		// ticker
		SaveTicker(ctx context.Context, ticker *Ticker) error
		FindTickers(ctx context.Context, assetID string) ([]*Ticker, error)
	}
)
