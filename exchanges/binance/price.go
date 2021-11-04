package binance

import (
	"context"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type (
	Price struct {
		Symbol string          `json:"symbol,omitempty"`
		Price  decimal.Decimal `json:"price,omitempty"`
	}
)

func (exch *binanceEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	resp, err := Request(ctx).SetQueryParam("symbol", symbol).Get("/ticker/bookTicker")
	if err != nil {
		log.WithError(err).Errorln("GET /ticker/price")
		return decimal.Zero, err
	}

	var ticker struct {
		Symbol   string          `json:"symbol"`
		BidPrice decimal.Decimal `json:"bidPrice"`
		BidQty   decimal.Decimal `json:"bidQty"`
		AskPrice decimal.Decimal `json:"askPrice"`
		AskQty   decimal.Decimal `json:"askQty"`
	}
	if err := UnmarshalResponse(resp, &ticker); err != nil {
		log.WithError(err).Errorln("getPrices.UnmarshalResponse")
		return decimal.Zero, err
	}

	return ticker.BidPrice, nil
}
