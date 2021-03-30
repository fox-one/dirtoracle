package binance

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
	exchangeName = "binance"
)

type binanceEx struct {
	cache *cache.Cache
}

func New() core.Exchange {
	return &binanceEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (b *binanceEx) Name() string {
	return exchangeName
}

func (b *binanceEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	pairSymbol := b.pairSymbol(b.assetSymbol(a.Symbol))
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"exchange": b.Name(),
		"symbol":   a.Symbol,
		"pair":     pairSymbol,
	})
	ctx = logger.WithContext(ctx, log)

	if ok, err := b.supported(ctx, pairSymbol); err != nil || !ok {
		return decimal.Zero, err
	}

	prices, err := b.getPrices(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	for _, price := range prices {
		if price.Symbol == pairSymbol {
			return price.Price, nil
		}
	}
	return decimal.Zero, nil
}

func (b *binanceEx) assetSymbol(symbol string) string {
	return symbol
}

func (b *binanceEx) pairSymbol(symbol string) string {
	return symbol + "BUSD"
}
