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
)

func provideAllExchanges() map[string]core.Exchange {
	block := exchanges.Block("XIN", "BOX")
	binance := block(provideBinanceExchanges())
	coinbase := block(provideCoinbaseExchanges())
	bitstamp := block(provideBitstampExchanges())
	bittrex := block(provideBittrexExchanges())
	bitfinex := block(provideBitfinexExchanges())
	huobi := block(provideHuobixchanges())

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

func provideHuobixchanges() core.Exchange {
	return huobi.New()
}
