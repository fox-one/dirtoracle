package bitstamp

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/fox-one/dirtoracle/core"
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
		exch = New()
		ctx  = context.Background()
	)

	for _, a := range assets {
		t.Run(exch.Name()+"-"+a.Symbol, func(t *testing.T) {
			p, err := exch.GetPrice(ctx, a)
			require.Nil(t, err, "GetPrice")
			t.Log(a.Symbol, "price:", p)
			require.True(t, p.IsPositive(), a.Symbol+" price not positive")
		})
	}

	{
		asset := &core.Asset{
			Symbol:  "XIN",
			AssetID: "c94ac88f-4671-3976-b60a-09064f1811e8",
		}
		p, err := exch.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("XIN price:", p)
		require.True(t, p.IsZero(), "XIN was not listed")
	}

	{
		asset := &core.Asset{
			Symbol:  "BOX",
			AssetID: "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
		}
		p, err := exch.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("BOX price:", p)
		require.True(t, p.IsZero(), "BOX was not listed")
	}
}
