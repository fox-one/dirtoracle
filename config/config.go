package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/pandodao/blst"
	"github.com/shopspring/decimal"
)

type (
	Config struct {
		DB    db.Config `json:"db"`
		Dapp  Dapp      `json:"dapp"`
		Group Group     `json:"group,omitempty"`
		Gas   Gas       `json:"gas,omitempty"`
	}

	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret"`
		Pin          string `json:"pin"`
	}

	Gas struct {
		Asset  string          `json:"asset,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`
	}

	Group struct {
		// 节点用于签名的私钥
		SignKey        *blst.PrivateKey `json:"sign_key,omitempty"`
		ConversationID string           `json:"conversation_id,omitempty"`
		Threshold      uint8            `json:"threshold,omitempty"`
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
