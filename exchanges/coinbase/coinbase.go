package coinbase

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "coinbase"
)

type coinbaseEx struct {
	cache *cache.Cache
}

func New() core.Exchange {
	return &coinbaseEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (*coinbaseEx) Name() string {
	return exchangeName
}

func (c *coinbaseEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	// block specific asset price from this exchange,
	//	since some assets were only be listed on 4swap,
	//	should avoid same symbol assets
	if c.IsAssetBlocked(ctx, a) {
		return decimal.Zero, nil
	}

	pairSymbol := c.pairSymbol(c.assetSymbol(a.Symbol))
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"exchange": c.Name(),
		"symbol":   a.Symbol,
		"pair":     pairSymbol,
	})
	ctx = logger.WithContext(ctx, log)

	if ok, err := c.supported(ctx, pairSymbol); err != nil || !ok {
		return decimal.Zero, err
	}

	ticker, err := c.getTicker(ctx, pairSymbol)
	if err != nil {
		return decimal.Zero, err
	}

	return ticker.Price, nil
}

func (*coinbaseEx) assetSymbol(symbol string) string {
	return symbol
}

func (*coinbaseEx) pairSymbol(symbol string) string {
	return symbol + "-USD"
}
