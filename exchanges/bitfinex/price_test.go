package bitfinex

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

var assets []*core.Asset

func init() {
	if bts, err := os.ReadFile("../testdata/assets.json"); err == nil {
		json.Unmarshal(bts, &assets)
	}
}

func TestGetPrice(t *testing.T) {
	var (
		exch = exchanges.Humanize(
			exchanges.PusdConverter(
				New(),
				fswap.New(),
				&core.Asset{
					AssetID: "9b180ab6-6abe-3dc0-a13f-04169eb34bfa",
					Symbol:  "USDC",
				},
				exchanges.PriceLimits{
					Min: decimal.New(90, -2),
					Max: decimal.New(110, -2),
				},
			))
		ctx = context.Background()
	)

	for _, a := range assets {
		t.Run(exch.Name()+"-"+a.Symbol, func(t *testing.T) {
			p, err := exch.GetPrice(ctx, a)
			t.Log(exch.Name(), a.Symbol, "price:", p)
			require.Nil(t, err, "GetPrice")
		})
	}
}
