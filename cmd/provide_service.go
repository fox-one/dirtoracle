package cmd

import (
	"github.com/fox-one/dirtoracle/core"
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

func provideWalletService(c *mixin.Client) core.WalletService {
	return wallet.New(c, cfg.Dapp.Pin)
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
