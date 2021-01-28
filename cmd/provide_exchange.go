package cmd

import (
	"github.com/fox-one/dirtoracle/core/exchange"
	"github.com/fox-one/dirtoracle/exchanges/binance"
	"github.com/fox-one/dirtoracle/exchanges/coinbase"
	"github.com/fox-one/dirtoracle/exchanges/exinswap"
	"github.com/fox-one/dirtoracle/exchanges/fswap"
)

func provideAllExchanges() map[string]exchange.Interface {
	fswap := provideFswapExchanges()
	eswap := provideExinswapExchanges()
	binance := provideBinanceExchanges()
	coinbase := provideCoinbaseExchanges()
	return map[string]exchange.Interface{
		fswap.Name():    fswap,
		eswap.Name():    eswap,
		binance.Name():  binance,
		coinbase.Name(): coinbase,
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
