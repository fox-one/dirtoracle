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

		price, err := c.Exchange.GetPrice(ctx, asset)
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

func FillSymbol(ex core.Exchange, assets core.AssetService) core.Exchange {
	return &assetEx{
		Exchange: ex,
		assets:   assets,
	}
}

type assetEx struct {
	core.Exchange
	assets core.AssetService
}

func (a *assetEx) GetPrice(ctx context.Context, asset *core.Asset) (decimal.Decimal, error) {
	if asset.Symbol == "" {
		a, err := a.assets.ReadAsset(ctx, asset.AssetID)
		if err != nil {
			if err == core.ErrAssetNotExist {
				return decimal.Zero, nil
			}
			return decimal.Zero, err
		}
		asset.Symbol = a.Symbol
	}

	return a.Exchange.GetPrice(ctx, asset)
}

func Humanize(ex core.Exchange) core.Exchange {
	return &humanizeEx{
		Exchange: ex,
	}
}

type humanizeEx struct {
	core.Exchange
}

func (c *humanizeEx) GetPrice(ctx context.Context, asset *core.Asset) (decimal.Decimal, error) {
	price, err := c.Exchange.GetPrice(ctx, asset)
	if err != nil {
		return price, err
	}

	return c.humanizeDecimal(price), nil
}

func (c *humanizeEx) humanizeDecimal(price decimal.Decimal) decimal.Decimal {
	for i := 0; i < 8; i++ {
		if price.GreaterThanOrEqual(decimal.New(1, 4-int32(i))) {
			return price.Truncate(int32(i))
		}
	}
	return price.Truncate(8)
}
