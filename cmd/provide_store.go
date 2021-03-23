package cmd

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/store/market"
	"github.com/fox-one/dirtoracle/store/pricedata"
	"github.com/fox-one/dirtoracle/store/subscriber"
	"github.com/fox-one/dirtoracle/store/wallet"
	"github.com/fox-one/pkg/property"
	"github.com/fox-one/pkg/store/db"
	propertystore "github.com/fox-one/pkg/store/property"
)

func provideDatabase() (*db.DB, error) {
	database, err := db.Open(cfg.DB)
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(database); err != nil {
		return nil, err
	}

	return database, nil
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

func provideSubscriberStore(db *db.DB) core.SubscriberStore {
	return subscriber.New(db)
}

func providePriceDataStore(db *db.DB) core.PriceDataStore {
	return pricedata.New(db)
}
