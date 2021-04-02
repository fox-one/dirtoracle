package binance

import (
	"context"
	"time"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

const (
	pricesKey = "prices"
)

type (
	Price struct {
		Symbol string          `json:"symbol,omitempty"`
		Price  decimal.Decimal `json:"price,omitempty"`
	}
)

func (b *binanceEx) getPrices(ctx context.Context) ([]*Price, error) {
	if prices, ok := b.cache.Get(pricesKey); ok {
		return prices.([]*Price), nil
	}

	log := logger.FromContext(ctx)
	resp, err := Request(ctx).Get("/ticker/price")
	if err != nil {
		log.WithError(err).Errorln("GET /ticker/price")
		return nil, err
	}

	var prices []*Price
	if err := UnmarshalResponse(resp, &prices); err != nil {
		log.WithError(err).Errorln("getPrices.UnmarshalResponse")
		return nil, err
	}

	b.cache.Set(pricesKey, prices, time.Second*10)
	return prices, nil
}
