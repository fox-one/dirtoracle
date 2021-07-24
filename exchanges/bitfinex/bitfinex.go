package bitfinex

import (
	"context"
	"fmt"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

const (
	exchangeName = "bitfinex"
)

type bitfinexEx struct{}

func New() core.Exchange {
	return &bitfinexEx{}
}

func (b *bitfinexEx) Name() string {
	return exchangeName
}

func (b *bitfinexEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	pairSymbol := b.pairSymbol(b.assetSymbol(a.Symbol))
	log := logger.FromContext(ctx).WithFields(logrus.Fields{
		"exchange": b.Name(),
		"symbol":   a.Symbol,
		"pair":     pairSymbol,
	})
	ctx = logger.WithContext(ctx, log)

	return b.getPrice(ctx, pairSymbol)
}

func (b *bitfinexEx) assetSymbol(symbol string) string {
	return symbol
}

func (b *bitfinexEx) pairSymbol(symbol string) string {
	switch symbol {
	case "BCH":
		return "tBCHN:USD"
	case "DOGE":
		return "tDOGE:USD"
	default:
		return fmt.Sprintf("t%sUSD", symbol)
	}
}
