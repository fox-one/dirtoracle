package ftx

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
	exchangeName = "ftx"
)

type ftxEx struct {
	cache *cache.Cache
}

func New() core.Exchange {
	return &ftxEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (*ftxEx) Name() string {
	return exchangeName
}

func (c *ftxEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
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

	return c.getPrice(ctx, pairSymbol)
}

func (*ftxEx) assetSymbol(symbol string) string {
	return symbol
}

func (*ftxEx) pairSymbol(symbol string) string {
	return symbol + "/USD"
}
