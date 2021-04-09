package exinswap

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
)

func (f *eswapEx) ListPortfolioTokens(ctx context.Context) ([]*core.PortfolioToken, error) {
	pairs, err := f.getPairs(ctx)
	if err != nil {
		return nil, err
	}

	var tokens []*core.PortfolioToken
	for _, pair := range pairs {
		if pair.Liquidity.IsZero() {
			continue
		}
		tokens = append(tokens, &core.PortfolioToken{
			AssetID: pair.LiquidityAssetID,
			Items: []*core.PortfolioItem{
				{
					AssetID: pair.BaseAssetID,
					Amount:  pair.BaseAmount.DivRound(pair.Liquidity, 8),
				},
				{
					AssetID: pair.QuoteAssetID,
					Amount:  pair.QuoteAmount.DivRound(pair.Liquidity, 8),
				},
			},
		})
	}

	return tokens, nil
}
