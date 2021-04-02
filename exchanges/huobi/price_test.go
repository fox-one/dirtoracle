package huobi

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
			Symbol: "BTC",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("BTC price:", p)
		require.True(t, p.IsPositive(), "BTC price not positive")
	}

	{
		asset := &core.Asset{
			Symbol: "XIN",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("XIN price:", p)
		require.True(t, p.IsZero(), "XIN was not listed")
	}

	{
		asset := &core.Asset{
			Symbol: "BOX",
		}
		p, err := b.GetPrice(ctx, asset)
		require.Nil(t, err, "GetPrice")
		t.Log("BOX price:", p)
		require.True(t, p.IsZero(), "BOX was not listed")
	}
}
