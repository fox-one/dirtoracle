package cmd

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/store/feeder"
	"github.com/fox-one/dirtoracle/store/market"
	"github.com/fox-one/dirtoracle/store/wallet"
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/store/db"
	propertystore "github.com/fox-one/pkg/store/property"
)

func provideDatabase() *db.DB {
	return db.MustOpen(cfg.DB)
}

func providePropertyStore(db *db.DB) property.Store {
	return propertystore.New(db)
}

func provideMarketStore() core.MarketStore {
	return market.New()
}

func provideWalletStore(db *db.DB) core.WalletStore {
	return wallet.New(db)
}

func provideFeederStore(db *db.DB) core.FeederStore {
	return feeder.New(db)
}
