package bitfinex

import (
	"context"
	"fmt"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "bitfinex"
)

type bitfinexEx struct {
	*exchanges.Exchange
	cache *cache.Cache
}

func New() core.Exchange {
	return &bitfinexEx{
		Exchange: exchanges.New(),
		cache:    cache.New(time.Minute, time.Minute),
	}
}

func (b *bitfinexEx) Name() string {
	return exchangeName
}

func (b *bitfinexEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	// block specific asset price from this exchange,
	//	since some assets were only be listed on 4swap,
	//	should avoid same symbol assets
	if b.IsAssetBlocked(ctx, a) {
		return decimal.Zero, nil
	}

	pairSymbol := b.pairSymbol(b.assetSymbol(a.Symbol))
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"exchange": b.Name(),
		"symbol":   a.Symbol,
		"pair":     pairSymbol,
	})
	ctx = logger.WithContext(ctx, log)

	prices, err := b.getPrices(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	if price, ok := prices[pairSymbol]; ok {
		return price, nil
	}
	return decimal.Zero, nil
}

func (b *bitfinexEx) assetSymbol(symbol string) string {
	return symbol
}

func (b *bitfinexEx) pairSymbol(symbol string) string {
	return fmt.Sprintf("t%sUSD", symbol)
}
