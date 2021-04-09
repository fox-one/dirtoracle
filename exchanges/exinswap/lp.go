package exinswap

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
)

func (f *eswapEx) UnpackAsset(ctx context.Context, asset *core.Asset) ([]*core.Portfolio, error) {
	pairs, err := f.getPairs(ctx)
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		if pair.LiquidityAssetID == asset.AssetID {
			if pair.Liquidity.IsZero() {
				return nil, nil
			}

			return []*core.Portfolio{
				{
					Asset: core.Asset{
						AssetID: pair.BaseAssetID,
					},
					Amount: pair.BaseAmount.DivRound(pair.Liquidity, 8),
				},
				{
					Asset: core.Asset{
						AssetID: pair.QuoteAssetID,
					},
					Amount: pair.QuoteAmount.DivRound(pair.Liquidity, 8),
				},
			}, nil
		}
	}

	return nil, nil
}
