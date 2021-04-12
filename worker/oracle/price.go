package oracle

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/shopspring/decimal"
)

func (m *Oracle) getPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	p, err := m.exchange.GetPrice(ctx, a)
	if err == core.ErrAssetNotExist {
		return decimal.Zero, nil
	}
	return p, err
}
