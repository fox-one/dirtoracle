package huobi

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
	exchangeName = "huobi"
)

type huobiEx struct {
	cache *cache.Cache
}

func New() core.Exchange {
	return &huobiEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (b *huobiEx) Name() string {
	return exchangeName
}

func (b *huobiEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
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

	return b.getPrice(ctx, pairSymbol)
}

func (b *huobiEx) assetSymbol(symbol string) string {
	return symbol
}

func (b *huobiEx) pairSymbol(symbol string) string {
	return strings.ToLower(symbol) + "husd"
}
