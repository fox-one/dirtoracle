package bittrex

import (
	"context"
	"testing"

	"github.com/fox-one/dirtoracle/pkg/route"
	"github.com/stretchr/testify/require"
)

func TestRoutes(t *testing.T) {
	var (
		exch = New().(*bittrexEx)
		ctx  = context.Background()
	)

	pairs, err := exch.getPairs(ctx)
	require.Nil(t, err, "getPairs")

	for _, a := range assets {
		t.Run(exch.Name()+"-"+a.Symbol, func(t *testing.T) {
			symbol := exch.assetSymbol(a.Symbol)
			routes, ok := route.FindRoutes(pairs, symbol, QuoteSymbol)
			t.Log(exch.Name(), a.Symbol, "routes:", ok, routes)
		})
	}
}
