package exinswap

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
	pusdAsset    = "31d2ea9c-95eb-3355-b65b-ba096853bc18"
	exchangeName = "exinswap"
)

var (
	pusdFunds = decimal.New(1, 3)
)

type (
	eswapEx struct {
		cache *cache.Cache
	}
)

func New() core.Exchange {
	return &eswapEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (*eswapEx) Name() string {
	return exchangeName
}

func (f *eswapEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"exchange": f.Name(),
		"symbol":   a.Symbol,
	})
	ctx = logger.WithContext(ctx, log)

	pairs, err := f.getPairs(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	return f.getPrice(ctx, a, pairs)
}
