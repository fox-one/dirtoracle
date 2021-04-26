package bittrex

import (
	"context"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "bittrex"
)

type bittrexEx struct{}

func New() core.Exchange {
	return &bittrexEx{}
}

func (*bittrexEx) Name() string {
	return exchangeName
}

func (b *bittrexEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
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
