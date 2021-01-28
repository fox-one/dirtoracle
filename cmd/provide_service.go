package cmd

import (
	"fmt"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/pkg/blst"
	"github.com/fox-one/mixin-sdk-go"
)

func provideMixinClient() *mixin.Client {
	c, err := mixin.NewFromKeystore(&cfg.Dapp.Keystore)
	if err != nil {
		panic(err)
	}

	return c
}

func provideSystem() *core.System {
	members := make([]*core.Member, 0, len(cfg.Group.Members))
	for _, m := range cfg.Group.Members {
		verifyKey, err := blst.DecodePublicKey(m.VerifyKey)
		if err != nil {
			panic(fmt.Errorf("decode verify key for member %s failed", m.ClientID))
		}

		members = append(members, &core.Member{
			ClientID:  m.ClientID,
			VerifyKey: verifyKey,
		})
	}

	signKey, err := blst.DecodePrivateKey(cfg.Group.SignKey)
	if err != nil {
		panic(fmt.Errorf("decode sign key failed"))
	}

	return &core.System{
		Admins:    cfg.Group.Admins,
		ClientID:  cfg.Dapp.ClientID,
		Members:   members,
		Threshold: cfg.Group.Threshold,
		SignKey:   signKey,
	}
}
