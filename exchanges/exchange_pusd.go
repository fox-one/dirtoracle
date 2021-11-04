package exchanges

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/singleflight"
)

type pusdEx struct {
	core.Exchange
	quoteAsset *core.Asset
	fswap      core.Exchange
	sf         *singleflight.Group
	cache      *cache.Cache
}

func (exch *pusdEx) GetPrice(ctx context.Context, asset *core.Asset) (decimal.Decimal, error) {
	price, err := exch.Exchange.GetPrice(ctx, asset)
	if err != nil || !price.IsPositive() {
		return price, err
	}

	convertRate, err := exch.GetConvertRate(ctx)
	if err != nil || !convertRate.IsPositive() {
		return decimal.Zero, err
	}

	return price.Mul(convertRate).Truncate(16), nil
}

func (exch *pusdEx) GetConvertRate(ctx context.Context) (decimal.Decimal, error) {
	key := "convert-tate"
	v, err, _ := exch.sf.Do(key, func() (interface{}, error) {
		if p, ok := exch.cache.Get(key); ok {
			return p, nil
		}

		quotePrice, err := exch.Exchange.GetPrice(ctx, exch.quoteAsset)
		if err != nil || !quotePrice.IsPositive() {
			return decimal.Zero, err
		}

		qQuotePrice, err := exch.fswap.GetPrice(ctx, exch.quoteAsset)
		if err != nil || !qQuotePrice.IsPositive() {
			return decimal.Zero, err
		}

		rate := qQuotePrice.Div(quotePrice)
		exch.cache.SetDefault(key, rate)
		return rate, nil
	})

	if err != nil {
		return decimal.Zero, err
	}

	return v.(decimal.Decimal), nil
}

func PusdConverter(exch, fswap core.Exchange, quoteAsset *core.Asset) core.Exchange {
	return &pusdEx{
		Exchange:   exch,
		fswap:      fswap,
		quoteAsset: quoteAsset,
		sf:         &singleflight.Group{},
		cache:      cache.New(time.Minute*10, time.Minute*10),
	}
}
