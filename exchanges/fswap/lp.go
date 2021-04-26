package fswap

import (
	"context"
	"time"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/dirtoracle/core"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
)

const (
	lpExchangeName = "4swap-lp"
)

type (
	lpEx struct {
		fswapEx
		core.Exchange
	}
)

func Lp(ex core.Exchange) core.Exchange {
	return &lpEx{
		fswapEx: fswapEx{
			cache: cache.New(time.Minute, time.Minute),
		},
		Exchange: ex,
	}
}

func (*lpEx) Name() string {
	return lpExchangeName
}

func (lp *lpEx) findLP(ctx context.Context, a *core.Asset) (*fswapsdk.Pair, error) {
	pairs, err := lp.getPairs(ctx)
	if err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		if pair.LiquidityAssetID == a.AssetID {
			return pair, nil
		}
	}
	return nil, nil
}
func (lp *lpEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	pair, err := lp.findLP(ctx, a)
	if err != nil {
		return decimal.Zero, err
	}

	if pair == nil {
		return lp.Exchange.GetPrice(ctx, a)
	}

	if pair.Liquidity.IsZero() {
		return decimal.Zero, nil
	}

	var assets = []struct {
		AssetID string
		Amount  decimal.Decimal
	}{
		{AssetID: pair.BaseAssetID, Amount: pair.BaseAmount.Div(pair.Liquidity)},
		{AssetID: pair.QuoteAssetID, Amount: pair.QuoteAmount.Div(pair.Liquidity)},
	}

	value := decimal.Zero
	for _, item := range assets {
		p, err := lp.Exchange.GetPrice(ctx, &core.Asset{AssetID: item.AssetID})
		if err != nil || p.IsZero() {
			return decimal.Zero, err
		}
		value = value.Add(p.Mul(item.Amount))
	}

	return value.Truncate(12), nil
}
