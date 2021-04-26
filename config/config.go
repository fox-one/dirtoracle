package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/pandodao/blst"
	"github.com/shopspring/decimal"
)

type (
	Config struct {
		DB     db.Config `json:"db"`
		Bwatch Bwatch    `json:"bwatch"`
		Dapp   Dapp      `json:"dapp"`
		Group  Group     `json:"group"`
		Gas    Gas       `json:"gas"`
	}

	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret"`
		Pin          string `json:"pin"`
	}

	Bwatch struct {
		ApiBase string `json:"api_base"`
	}

	Gas struct {
		Asset  string          `json:"asset"`
		Amount decimal.Decimal `json:"amount"`
	}

	Group struct {
		// 节点用于签名的私钥
		SignKey        *blst.PrivateKey `json:"sign_key"`
		ConversationID string           `json:"conversation_id"`
		Threshold      uint8            `json:"threshold"`
	}
)

func defaultVote(cfg *Config) {
	if cfg.Gas.Asset == "" {
		cfg.Gas.Asset = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	}

	if cfg.Gas.Amount.IsZero() {
		cfg.Gas.Amount = decimal.New(1, -8)
	}
}
