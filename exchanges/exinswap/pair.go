package exinswap

import (
	"context"
	"time"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

const (
	pairsKey = "pairs"
)

func (c *eswapEx) getPairs(ctx context.Context) ([]*fswapsdk.Pair, error) {
	if pairs, ok := c.cache.Get(pairsKey); ok {
		return pairs.([]*fswapsdk.Pair), nil
	}

	log := logger.FromContext(ctx)
	const uri = "/pairs"
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return nil, err
	}

	var body struct {
		Pairs []struct {
			Asset0 struct {
				ID    string          `json:"uuid"`
				Price decimal.Decimal `json:"priceUsdt"`
			} `json:"asset0"`
			Asset1 struct {
				ID    string          `json:"uuid"`
				Price decimal.Decimal `json:"priceUsdt"`
			} `json:"asset1"`
			LPAsset struct {
				ID    string          `json:"uuid"`
				Price decimal.Decimal `json:"priceUsdt"`
			} `json:"lpAsset"`
			Balance0      decimal.Decimal `json:"asset0Balance"`
			Balance1      decimal.Decimal `json:"asset1Balance"`
			LPAssetSupply decimal.Decimal `json:"lpAssetSupply"`
			Volume        decimal.Decimal `json:"usdtTradeVolume24hours"`
			TradeType     string          `json:"tradeType"`
		} `json:"data"`
	}
	if err := UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("getPairs.UnmarshalResponse")
		return nil, err
	}

	var (
		pairs    = make([]*fswapsdk.Pair, len(body.Pairs))
		fee      = decimal.New(3, -3)
		curveFee = decimal.New(4, -4)
	)

	for i, p := range body.Pairs {
		pair := &fswapsdk.Pair{
			RouteID:          int64(i),
			BaseAssetID:      p.Asset0.ID,
			BaseAmount:       p.Balance0,
			QuoteAssetID:     p.Asset1.ID,
			QuoteAmount:      p.Balance1,
			LiquidityAssetID: p.LPAsset.ID,
			Liquidity:        p.LPAssetSupply,
			FeePercent:       fee,
		}

		if p.TradeType == "curve" {
			pair.SwapMethod = "curve"
			pair.FeePercent = curveFee
		}

		pairs[i] = pair
	}

	c.cache.Set(pairsKey, pairs, time.Second*10)
	return pairs, nil
}
