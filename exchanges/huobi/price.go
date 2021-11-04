package huobi

import (
	"context"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Symbol string  `json:"symbol"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
		Volume float64 `json:"vol"`
	}
)

func (b *huobiEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	uri := "/market/detail/merged"
	resp, err := Request(ctx).SetQueryParam("symbol", symbol).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return decimal.Zero, err
	}

	var body struct {
		Tick struct {
			ID     uint64     `json:"id"`
			Open   float64    `json:"open"`
			Close  float64    `json:"close"`
			Low    float64    `json:"low"`
			High   float64    `json:"high"`
			Volume float64    `json:"vol"`
			Ask    [2]float64 `json:"ask"`
			Bid    [2]float64 `json:"bid"`
		} `json:"tick"`
	}
	if err := UnmarshalResponse(resp, &body); err != nil {
		log.WithError(err).Errorln("getTicker.UnmarshalResponse")
		return decimal.Zero, err
	}

	return decimal.NewFromFloat(body.Tick.Bid[0]), nil
}
