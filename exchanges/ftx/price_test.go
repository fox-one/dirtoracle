package ftx

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
	"github.com/stretchr/testify/require"
)

var assets []*core.Asset

func init() {
	if bts, err := ioutil.ReadFile("../testdata/assets.json"); err == nil {
		json.Unmarshal(bts, &assets)
	}
}

func TestGetPrice(t *testing.T) {
	var (
		exch = exchanges.Humanize(exchanges.PusdConverter(New(), fswap.New(), &core.Asset{
			AssetID: "8549b4ad-917c-3461-a646-481adc5d7f7f",
			Symbol:  "DAI",
		}))
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
