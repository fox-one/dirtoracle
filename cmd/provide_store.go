package cmd

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/store/subscriber"
	"github.com/fox-one/dirtoracle/store/wallet"
	"github.com/fox-one/pkg/store/db"
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

func provideWalletStore(db *db.DB) core.WalletStore {
	return wallet.New(db)
}

func provideSubscriberStore(db *db.DB) core.SubscriberStore {
	return subscriber.New(db)
}
