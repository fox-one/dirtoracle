package fswap

import (
	"context"
	"time"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	pusdAsset    = "31d2ea9c-95eb-3355-b65b-ba096853bc18"
	exchangeName = "4swap"
)

var (
	pusdFunds = decimal.New(1, 3)
)

type (
	fswapEx struct {
		cache *cache.Cache
	}
)

func init() {
	fswapsdk.UseEndpoint("https://lake-api.pando.im")
}

func New() core.Exchange {
	return &fswapEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (*fswapEx) Name() string {
	return exchangeName
}

func (f *fswapEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	if a.AssetID == pusdAsset {
		return decimal.New(1, 0), nil
	}

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
