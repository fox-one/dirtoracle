package cmd

import (
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/dirtoracle/exchanges/binance"
	"github.com/fox-one/dirtoracle/exchanges/bitstamp"
	"github.com/fox-one/dirtoracle/exchanges/bittrex"
	"github.com/fox-one/dirtoracle/exchanges/coinbase"
	"github.com/fox-one/dirtoracle/exchanges/exinswap"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
	"github.com/fox-one/dirtoracle/exchanges/kraken"
)

func provideAllExchanges() map[string]exchange.Interface {
	fswap := provideFswapExchanges()
	eswap := provideExinswapExchanges()
	binance := provideBinanceExchanges()
	coinbase := provideCoinbaseExchanges()
	bitstamp := provideBitstampExchanges()
	kraken := provideKrakenExchanges()
	bittrex := provideBittrexExchanges()
	return map[string]exchange.Interface{
		fswap.Name():    fswap,
		eswap.Name():    eswap,
		binance.Name():  binance,
		coinbase.Name(): coinbase,
		bitstamp.Name(): bitstamp,
		kraken.Name():   kraken,
		bittrex.Name():  bittrex,
	}
}

func provideFswapExchanges() exchange.Interface {
	return fswap.New()
}

func provideExinswapExchanges() exchange.Interface {
	return exinswap.New()
}

func provideBinanceExchanges() exchange.Interface {
	return binance.New()
}

func provideCoinbaseExchanges() exchange.Interface {
	return coinbase.New()
}

func provideBitstampExchanges() exchange.Interface {
	return bitstamp.New()
}

func provideKrakenExchanges() exchange.Interface {
	return kraken.New()
}

func provideBittrexExchanges() exchange.Interface {
	return bittrex.New()
}
