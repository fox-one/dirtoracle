package bitfinex

import (
	"context"
	"testing"

	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/stretchr/testify/require"
)

func TestRoutes(t *testing.T) {
	var (
		exch = New().(*bitfinexEx)
		ctx  = context.Background()
	)

	pairs, err := exch.getPairs(ctx)
	require.Nil(t, err, "getPairs")

	for _, a := range assets {
		t.Run(exch.Name()+"-"+a.Symbol, func(t *testing.T) {
			symbol := exch.assetSymbol(a.Symbol)
			routes, ok := route.FindRoutes(pairs, symbol, QuoteSymbol)
			t.Log(exch.Name(), a.Symbol, "routes:", routes)
			require.True(t, ok, "FindRoutes")
			require.NotEmpty(t, routes, "empty routes")
		})
	}
}
