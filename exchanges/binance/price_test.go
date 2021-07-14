package binance

import (
	"context"
	"testing"

	"github.com/fox-one/dirtoracle/core"
	"github.com/stretchr/testify/require"
)

var assets = []*core.Asset{
	{
		Symbol:  "BTC",
		AssetID: "c6d0c728-2624-429b-8e0d-d9d19b6592fa",
	},
	{
		Symbol:  "ETH",
		AssetID: "43d61dcd-e413-450d-80b8-101d5e903357",
	},
	{
		Symbol:  "EOS",
		AssetID: "6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
	},
	{
		Symbol:  "DOGE",
		AssetID: "6770a1e5-6086-44d5-b60f-545f9d9e8ffd",
	},
	{
		Symbol:  "ZEC",
		AssetID: "c996abc9-d94e-4494-b1cf-2a3fd3ac5714",
	},
	{
		Symbol:  "DOT",
		AssetID: "54c61a72-b982-4034-a556-0d99e3c21e39",
	},
	{
		Symbol:  "LTC",
		AssetID: "76c802a2-7c88-447f-a93e-c29c9e5dd9c8",
	},
	{
		Symbol:  "SC",
		AssetID: "990c4c29-57e9-48f6-9819-7d986ea44985",
	},
	{
		Symbol:  "ZEN",
		AssetID: "a2c5d22b-62a2-4c13-b3f0-013290dbac60",
	},
	{
		Symbol:  "BCH",
		AssetID: "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0",
	},
	{
		Symbol:  "FIL",
		AssetID: "08285081-e1d8-4be6-9edc-e203afa932da",
	},
}

func TestGetPrice(t *testing.T) {
	var (
		b   = New()
		ctx = context.Background()
	)

	for _, a := range assets {
		t.Run(b.Name()+"-"+a.Symbol, func(t *testing.T) {
			p, err := b.GetPrice(ctx, a)
			require.Nil(t, err, "GetPrice")
			t.Log(a.Symbol, "price:", p)
			require.True(t, p.IsPositive(), "BTC price not positive")
		})

	}

	{
		asset := &core.Asset{
			Symbol:  "XIN",
			AssetID: "c94ac88f-4671-3976-b60a-09064f1811e8",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("XIN price:", p)
		require.True(t, p.IsZero(), "XIN was not listed")
	}

	{
		asset := &core.Asset{
			Symbol:  "BOX",
			AssetID: "f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("BOX price:", p)
		require.True(t, p.IsZero(), "BOX was not listed")
	}
}
