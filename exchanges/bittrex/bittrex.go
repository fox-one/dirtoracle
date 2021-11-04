package bittrex

import (
	"context"
	"time"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
)

const (
	exchangeName = "bittrex"

	QuoteSymbol = "USD"
)

type bittrexEx struct {
	cache *cache.Cache
}

func New() core.Exchange {
	return &bittrexEx{
		cache: cache.New(time.Minute, time.Minute),
	}
}

func (*bittrexEx) Name() string {
	return exchangeName
}

func (exch *bittrexEx) GetPrice(ctx context.Context, a *core.Asset) (decimal.Decimal, error) {
	symbol := exch.assetSymbol(a.Symbol)
	if symbol == QuoteSymbol {
		return decimal.New(1, 0), nil
	}

	pairs, err := exch.getPairs(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	routes, ok := route.FindRoutes(pairs, symbol, QuoteSymbol)
	if !ok {
		return decimal.Zero, err
	}

	var price = decimal.New(1, 0)
	for _, route := range routes {
		p, err := exch.getPrice(ctx, route.Symbol)
		if err != nil {
			return decimal.Zero, err
		}
		if route.Reverse {
			price = price.Div(p)
		} else {
			price = price.Mul(p)
		}
	}

	return price, nil
}

func (*bittrexEx) assetSymbol(symbol string) string {
	return symbol
}
