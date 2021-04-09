package oracle

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func (m *Oracle) getPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	{
		pos, err := m.unpackAsset(ctx, a.AssetID)
		if err != nil {
			return decimal.Zero, err
		}

		if len(pos) > 0 {
			var value decimal.Decimal
			for _, a := range pos {
				p, err := m.getPrice(ctx, &core.Asset{AssetID: a.AssetID})
				if err != nil || !p.IsPositive() {
					return decimal.Zero, err
				}
				value = value.Add(p.Mul(a.Amount))
			}
			return value.Truncate(8), nil
		}
	}

	if a.Symbol == "" {
		t, err := m.getAsset(ctx, a.AssetID)
		if err != nil {
			logger.FromContext(ctx).WithError(err).Errorln("read asset failed")
			return decimal.Zero, err
		}
		if t == nil {
			logger.FromContext(ctx).WithError(err).Errorln("asset not found")
			return decimal.Zero, nil
		}
		a.Symbol = t.Symbol
	}

	for _, e := range m.exchanges {
		p, err := e.GetPrice(ctx, a)
		if err != nil || p.IsPositive() {
			return p, err
		}
	}
	return decimal.Zero, nil
}
