package bwatch

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
)

const (
	serviceName = "bwatch"
)

type bwatchService struct {
	core.Exchange
	cache  *cache.Cache
	assets core.AssetService
}

func New(ex core.Exchange, assets core.AssetService) core.Exchange {
	return &bwatchService{
		Exchange: ex,
		assets:   assets,
		cache:    cache.New(time.Hour, time.Minute),
	}
}

func (bwatchService) Name() string {
	return serviceName
}

func (b *bwatchService) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	etf, err := b.getETF(ctx, a)
	if err != nil {
		return decimal.Zero, err
	}

	if etf == nil {
		return b.Exchange.GetPrice(ctx, a)
	}

	value := decimal.Zero
	for id, v := range etf.Assets {
		a, err := b.assets.ReadAsset(ctx, id)
		if err != nil {
			return decimal.Zero, err
		}
		p, err := b.Exchange.GetPrice(ctx, a)
		if err != nil {
			return decimal.Zero, err
		}
		value = value.Add(p.Mul(v))
	}
	return value.Truncate(12), nil
}
