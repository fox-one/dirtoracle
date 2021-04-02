package exinswap

import (
	"context"
	"testing"

	"github.com/fox-one/dirtoracle/core"
	"github.com/stretchr/testify/require"
)

func TestGetPrice(t *testing.T) {
	var (
		b   = New()
		ctx = context.Background()
	)

	{
		asset := &core.Asset{
			Symbol:  "BTC",
			AssetID: "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("BTC price:", p)
		require.True(t, p.IsPositive(), "BTC price not positive")
	}

	{
		asset := &core.Asset{
			Symbol:  "XIN",
			AssetID: "c94ac88f-4671-3976-b60a-09064f1811e8",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("XIN price:", p)
		require.True(t, p.IsPositive(), "XIN was not listed")
	}

	{
		asset := &core.Asset{
			Symbol:  "BOX",
			AssetID: "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("BOX price:", p)
		require.True(t, p.IsPositive(), "BOX was not listed")
	}
}
