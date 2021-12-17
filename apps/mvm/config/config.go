package config

import (
	"github.com/fox-one/dirtoracle/apps/mvm/core"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type (
	Config struct {
		DB     db.Config `json:"db"`
		Dapp   Dapp      `json:"dapp"`
		Oracle Oracle    `json:"oracle"`
		MVM    MVM       `json:"mvm,omitempty"`
		Gas    Gas       `json:"gas,omitempty"`
	}

	Oracle struct {
		Signers   []*core.Signer `json:"signers"`
		Threshold uint8          `json:"threshold,omitempty"`
	}

	Dapp struct {
		mixin.Keystore
		Pin string `json:"pin"`
	}

	MVM struct {
		Process   string   `json:"process"`
		Groups    []string `json:"groups,omitempty"`
		Threshold uint8    `json:"threshold,omitempty"`
	}

	Gas struct {
		Asset  uuid.UUID       `json:"asset,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`
	}
)
