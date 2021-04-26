package bitstamp

import (
	"context"
	"strings"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "bitstamp"
)

type bitstampEx struct {
	cache *cache.Cache
}

func New() core.Exchange {
	return &bitstampEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (*bitstampEx) Name() string {
	return exchangeName
}

func (b *bitstampEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
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

	ticker, err := b.getTicker(ctx, pairSymbol)
	if err != nil {
		return decimal.Zero, err
	}

	return ticker.Last, nil
}

func (*bitstampEx) assetSymbol(symbol string) string {
	return symbol
}

func (*bitstampEx) pairSymbol(symbol string) string {
	return strings.ToLower(symbol) + "usd"
}
