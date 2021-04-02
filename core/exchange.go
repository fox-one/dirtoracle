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
)
