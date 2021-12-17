package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type (
	Asset struct {
		ID             uint            `sql:"PRIMARY_KEY" json:"id"`
		AssetID        string          `json:"asset_id,omitempty"`
		Symbol         string          `json:"symbol,omitempty"`
		Price          decimal.Decimal `sql:"type:decimal(24,12)" json:"price,omitempty"`
		PriceDuration  int64           `json:"price_duration,omitempty"`
		PriceUpdatedAt *time.Time      `json:"price_updated_at,omitempty"`
	}

	AssetStore interface {
		List(ctx context.Context) ([]*Asset, error)
		Find(ctx context.Context, assetID string) (*Asset, error)
		Update(ctx context.Context, asset *Asset) error
	}
)
