package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/shopspring/decimal"
)

type (
	Config struct {
		DB    db.Config `json:"db"`
		Dapp  Dapp      `json:"dapp"`
		Group Group     `json:"group,omitempty"`
	}

	Redis struct {
		Addr string `json:"addr,omitempty"`
		DB   int    `json:"db,omitempty"`
	}

	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret"`
		Pin          string `json:"pin"`
	}

	Member struct {
		ClientID string `json:"client_id,omitempty"`
		// 节点共享的用户验证签名的公钥
		VerifyKey string `json:"verify_key,omitempty"`
	}

	Vote struct {
		Asset  string          `json:"asset,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`
	}

	Group struct {
		// 节点管理员 mixin id
		Admins []string `json:"admins,omitempty"`
		// 节点用于签名的私钥
		SignKey string `json:"sign_key,omitempty"`
		// 节点共享的用户解密的私钥
		PrivateKey string   `json:"private_key,omitempty"`
		Members    []Member `json:"members,omitempty"`
		Threshold  uint8    `json:"threshold,omitempty"`

		Vote Vote `json:"vote,omitempty"`
	}
)

func defaultVote(cfg *Config) {
	if cfg.Group.Vote.Asset == "" {
		cfg.Group.Vote.Asset = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	}

	if cfg.Group.Vote.Amount.IsZero() {
		cfg.Group.Vote.Amount = decimal.NewFromInt(1)
	}
}
