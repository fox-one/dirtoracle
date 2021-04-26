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
	"github.com/fox-one/dirtoracle/exchanges/huobi"
	"github.com/fox-one/dirtoracle/exchanges/okex"
)

func provideAllExchanges(assets core.AssetService) map[string]core.Exchange {
	block := exchanges.Block("f5ef6b5d-cc5a-3d90-b2c0-a2fd386e7a3c", "c94ac88f-4671-3976-b60a-09064f1811e8")
	binance := block(exchanges.FillSymbol(provideBinanceExchanges(), assets))
	coinbase := block(exchanges.FillSymbol(provideCoinbaseExchanges(), assets))
	bitstamp := block(exchanges.FillSymbol(provideBitstampExchanges(), assets))
	bittrex := block(exchanges.FillSymbol(provideBittrexExchanges(), assets))
	bitfinex := block(exchanges.FillSymbol(provideBitfinexExchanges(), assets))
	huobi := block(exchanges.FillSymbol(provideHuobiExchanges(), assets))
	okex := block(exchanges.FillSymbol(provideOkexExchanges(), assets))

	fswap := provideFswapExchanges()
	eswap := provideExinswapExchanges()

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
