package exchanges

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/dirtoracle/core"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/singleflight"
)

func Block(ids ...string) func(exchange core.Exchange) core.Exchange {
	return func(ex core.Exchange) core.Exchange {
		return &blockAssets{
			Exchange:  ex,
			blacklist: ids,
		}
	}
}

type blockAssets struct {
	core.Exchange
	blacklist []string
}

func (ex *blockAssets) GetPrice(ctx context.Context, asset *core.Asset) (decimal.Decimal, error) {
	if govalidator.IsIn(asset.AssetID, ex.blacklist...) {
		return decimal.Zero, nil
	}

	return ex.Exchange.GetPrice(ctx, asset)
}

func Chain(exs ...core.Exchange) core.Exchange {
	return chains(exs)
}

type chains []core.Exchange

func (c chains) Name() string {
	return "chains"
}

func (c chains) GetPrice(ctx context.Context, asset *core.Asset) (decimal.Decimal, error) {
	for _, ex := range c {
		p, err := ex.GetPrice(ctx, asset)
		if err != nil {
			return decimal.Zero, err
		}

		if p.IsZero() {
			continue
		}

		return p, nil
	}

	return decimal.Zero, nil
}

func Cache(ex core.Exchange, exp time.Duration) core.Exchange {
	return &cacheEx{
		Exchange: ex,
		cache:    cache.New(exp, time.Minute),
		sf:       &singleflight.Group{},
	}
}

type cacheEx struct {
	core.Exchange
	cache *cache.Cache
	sf    *singleflight.Group
}

func (c *cacheEx) GetPrice(ctx context.Context, asset *core.Asset) (decimal.Decimal, error) {
	v, err, _ := c.sf.Do(asset.AssetID, func() (interface{}, error) {
		if p, ok := c.cache.Get(asset.AssetID); ok {
			return p, nil
		}

		price, err := c.GetPrice(ctx, asset)
		if err != nil {
			return decimal.Zero, err
		}

		c.cache.SetDefault(asset.AssetID, price)
		return price, nil
	})

	if err != nil {
		return decimal.Zero, err
	}

	return v.(decimal.Decimal), nil
}
