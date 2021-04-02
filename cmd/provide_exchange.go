package cmd

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/exchanges/binance"
	"github.com/fox-one/dirtoracle/exchanges/bitstamp"
	"github.com/fox-one/dirtoracle/exchanges/bittrex"
	"github.com/fox-one/dirtoracle/exchanges/coinbase"
	"github.com/fox-one/dirtoracle/exchanges/exinswap"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
)

func provideAllExchanges() map[string]core.Exchange {
	fswap := provideFswapExchanges()
	eswap := provideExinswapExchanges()
	binance := provideBinanceExchanges()
	coinbase := provideCoinbaseExchanges()
	bitstamp := provideBitstampExchanges()
	bittrex := provideBittrexExchanges()
	return map[string]core.Exchange{
		fswap.Name():    fswap,
		eswap.Name():    eswap,
		binance.Name():  binance,
		coinbase.Name(): coinbase,
		bitstamp.Name(): bitstamp,
		bittrex.Name():  bittrex,
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
