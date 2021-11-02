package ftx

import (
	"context"
	"fmt"

	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func (exch *ftxEx) getPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)
	uri := fmt.Sprintf("/markets/%s", symbol)
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		log.WithError(err).Errorln("GET", uri)
		return decimal.Zero, err
	}

	var pair Pair
	if err := UnmarshalResponse(resp, &pair); err != nil {
		log.WithError(err).Errorln("getPrice.UnmarshalResponse")
		return decimal.Zero, err
	}

	return decimal.NewFromFloat(pair.Bid), nil
}
