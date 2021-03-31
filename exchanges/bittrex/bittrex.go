package bittrex

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges"
	"github.com/fox-one/pkg/logger"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "bittrex"
)

type bittrexEx struct {
	*exchanges.Exchange
	cache *cache.Cache
}

func New() core.Exchange {
	return &bittrexEx{
		Exchange: exchanges.New(),
		cache:    cache.New(time.Minute, time.Minute),
	}
}

func (*bittrexEx) Name() string {
	return exchangeName
}

func (b *bittrexEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
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

	tickers, err := b.getTickers(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	for _, ticker := range tickers {
		if ticker.Symbol == pairSymbol {
			return ticker.LastTradeRate, nil
		}
	}

	return decimal.Zero, nil
}

func (*bittrexEx) assetSymbol(symbol string) string {
	return symbol
}

func (*bittrexEx) pairSymbol(symbol string) string {
	return symbol + "-USD"
}
