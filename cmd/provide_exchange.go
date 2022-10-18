package cmd

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges"
	"github.com/fox-one/dirtoracle/exchanges/binance"
	"github.com/fox-one/dirtoracle/exchanges/bitfinex"
	"github.com/fox-one/dirtoracle/exchanges/bitstamp"
	"github.com/fox-one/dirtoracle/exchanges/bittrex"
	"github.com/fox-one/dirtoracle/exchanges/coinbase"
	"github.com/fox-one/dirtoracle/exchanges/exinswap"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
	"github.com/fox-one/dirtoracle/exchanges/ftx"
	"github.com/fox-one/dirtoracle/exchanges/huobi"
	"github.com/fox-one/dirtoracle/exchanges/okex"
)

func provideAllExchanges(assets core.AssetService) map[string]core.Exchange {
	fswap := provideFswapExchanges()
	eswap := provideExinswapExchanges()

	var (
		usdc = &core.Asset{
			AssetID: "9b180ab6-6abe-3dc0-a13f-04169eb34bfa",
			Symbol:  "USDC",
		}

		priceLimts = exchanges.PriceLimits{}
	)

	for _, limits := range cfg.PriceLimits {
		if limits.AssetID == usdc.AssetID {
			priceLimts.Max = limits.Max
			priceLimts.Min = limits.Min
		}
	}

	block := exchanges.Block("f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c", "c94ac88f-4671-3976-b60a-09064f1811e8")
	wrapExchangeFunc := func(exch core.Exchange) core.Exchange {
		return block(
			exchanges.PusdConverter(
				exchanges.FillSymbol(
					exch,
					assets,
				),
				fswap,
				usdc,
				priceLimts,
			),
		)
	}

	binance := wrapExchangeFunc(provideBinanceExchanges())
	coinbase := wrapExchangeFunc(provideCoinbaseExchanges())
	bitstamp := wrapExchangeFunc(provideBitstampExchanges())
	bittrex := wrapExchangeFunc(provideBittrexExchanges())
	bitfinex := wrapExchangeFunc(provideBitfinexExchanges())
	huobi := wrapExchangeFunc(provideHuobiExchanges())
	okex := wrapExchangeFunc(provideOkexExchanges())
	ftx := wrapExchangeFunc(provideFtxExchanges())

	return map[string]core.Exchange{
		fswap.Name():    fswap,
		eswap.Name():    eswap,
		binance.Name():  binance,
		coinbase.Name(): coinbase,
		bitstamp.Name(): bitstamp,
		bittrex.Name():  bittrex,
		bitfinex.Name(): bitfinex,
		huobi.Name():    huobi,
		okex.Name():     okex,
		ftx.Name():      ftx,
	}
}

func provideFswapExchanges() core.Exchange {
	return fswap.New()
}

func provideExinswapExchanges() core.Exchange {
	return exinswap.New()
}

func provideBinanceExchanges() core.Exchange {
	return binance.New()
}

func provideCoinbaseExchanges() core.Exchange {
	return coinbase.New()
}

func provideBitstampExchanges() core.Exchange {
	return bitstamp.New()
}

func provideBittrexExchanges() core.Exchange {
	return bittrex.New()
}

func provideBitfinexExchanges() core.Exchange {
	return bitfinex.New()
}

func provideHuobiExchanges() core.Exchange {
	return huobi.New()
}

func provideOkexExchanges() core.Exchange {
	return okex.New()
}

func provideFtxExchanges() core.Exchange {
	return ftx.New()
}
