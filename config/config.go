package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/pandodao/blst"
	"github.com/pandodao/blst/en256"
	"github.com/shopspring/decimal"
)

type (
	Config struct {
		DB          db.Config      `json:"db"`
		Bwatch      Bwatch         `json:"bwatch"`
		Dapp        Dapp           `json:"dapp"`
		Group       Group          `json:"group"`
		Gas         Gas            `json:"gas"`
		PriceLimits []*PriceLimits `json:"price_limits"`
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
		SignKey        *blst.PrivateKey  `json:"sign_key"`
		En256SignKey   *en256.PrivateKey `json:"en256_sign_key"`
		ConversationID string            `json:"conversation_id"`
	}

	PriceLimits struct {
		AssetID string          `json:"asset_id"`
		Min     decimal.Decimal `json:"min"`
		Max     decimal.Decimal `json:"max"`
	}
)
