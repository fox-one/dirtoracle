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

	Portfolio struct {
		Asset
		Amount decimal.Decimal `json:"amount"`
	}

	PortfolioService interface {
		UnpackAsset(ctx context.Context, asset *Asset) ([]*Portfolio, error)
	}
)
