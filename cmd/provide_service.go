package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/dirtoracle/service/wallet"
	"github.com/fox-one/mixin-sdk-go"
	homedir "github.com/mitchellh/go-homedir"
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
		Admins:         cfg.Group.Admins,
		ClientID:       cfg.Dapp.ClientID,
		Members:        cfg.Group.Members,
		Threshold:      cfg.Group.Threshold,
		SignKey:        cfg.Group.SignKey,
		ConversationID: cfg.Group.ConversationID,
		GasAsset:       cfg.Gas.Asset,
		GasAmount:      cfg.Gas.Amount,
	}

	if s.Me() == nil {
		panic(fmt.Errorf("dapp is not a group member"))
	}

	d := map[int64]bool{}
	for _, m := range s.Members {
		if m.ID < 0 || m.ID >= 64 {
			panic(fmt.Errorf("invalid: group member id (%d)", m.ID))
		}
		if _, ok := d[m.ID]; ok {
			panic(fmt.Errorf("repeated group member id (%d)", m.ID))
		}
		d[m.ID] = true
	}

	return s
}

func provideFeedConfigs(feedFile string) []*core.FeedConfig {
	if feedFile == "" {
		dir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}

		filename := path.Join(dir, ".feeds.json")
		info, err := os.Stat(filename)
		if !os.IsNotExist(err) && !info.IsDir() {
			feedFile = filename
		}
	}

	bts, err := ioutil.ReadFile(feedFile)
	if err != nil {
		panic(err)
	}

	var feeds []*core.FeedConfig
	if err := json.Unmarshal(bts, &feeds); err != nil {
		panic(err)
	}
	return feeds
}
