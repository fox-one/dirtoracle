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
	cache *cache.Cache
}

func New(ex core.Exchange) core.Exchange {
	return &bwatchService{
		Exchange: ex,
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
		p, err := b.Exchange.GetPrice(ctx, &core.Asset{AssetID: id})
		if err != nil {
			return decimal.Zero, err
		}
		value = value.Add(p.Mul(v))
	}
	return value.Truncate(12), nil
}
