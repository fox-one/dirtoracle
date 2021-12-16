package cmd

import "github.com/fox-one/mixin-sdk-go"

func provideMixinClient() *mixin.Client {
	c, err := mixin.NewFromKeystore(&cfg.Dapp)
	if err != nil {
		panic(err)
	}

	return c
}
