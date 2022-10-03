package exinswap

import (
	"context"

	fswapsdk "github.com/fox-one/4swap-sdk-go"
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/pkg/logger"
	"github.com/shopspring/decimal"
)

func (f *eswapEx) getPrice(ctx context.Context, asset *core.Asset, pairs []*fswapsdk.Pair) (decimal.Decimal, error) {
	log := logger.FromContext(ctx)

	order, err := fswapsdk.Route(pairs, pusdAsset, asset.AssetID, pusdFunds)
	if err != nil {
		log.WithError(err).Errorln("Route")
		return decimal.Zero, nil
	}

	return order.PayAmount.Div(order.FillAmount).Truncate(8), nil
}
