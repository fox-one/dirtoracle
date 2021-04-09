package core

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	Exchange interface {
		Name() string
		GetPrice(ctx context.Context, asset *Asset) (decimal.Decimal, error)
	}

	PortfolioItem struct {
		AssetID string          `json:"asset_id"`
		Amount  decimal.Decimal `json:"amount"`
	}

	PortfolioToken struct {
		AssetID string           `json:"asset_id"`
		Items   []*PortfolioItem `json:"items"`
	}

	PortfolioService interface {
		Name() string
		ListPortfolioTokens(ctx context.Context) ([]*PortfolioToken, error)
	}
)
