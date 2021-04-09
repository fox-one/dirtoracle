package cmd

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/service/asset"
	"github.com/fox-one/dirtoracle/service/bwatch"
	"github.com/fox-one/dirtoracle/service/wallet"
	"github.com/fox-one/mixin-sdk-go"
)

func provideMixinClient() *mixin.Client {
	c, err := mixin.NewFromKeystore(&cfg.Dapp.Keystore)
	if err != nil {
		panic(err)
	}

	return c
}

func provideAssetService(c *mixin.Client) core.AssetService {
	return asset.New(c)
}

func provideWalletService(c *mixin.Client) core.WalletService {
	return wallet.New(c, cfg.Dapp.Pin)
}

func provideBwatchService() core.PortfolioService {
	return bwatch.New(bwatch.Config{ApiHost: cfg.Bwatch.ApiBase})
}

func provideSystem() *core.System {
	s := &core.System{
		ConversationID: cfg.Group.ConversationID,
		SignKey:        cfg.Group.SignKey,
		VerifyKey:      cfg.Group.SignKey.PublicKey(),
		GasAsset:       cfg.Gas.Asset,
		GasAmount:      cfg.Gas.Amount,
	}

	return s
}
