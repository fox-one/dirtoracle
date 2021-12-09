package config

import (
	"github.com/fox-one/dirtoracle/core"
	"github.com/fox-one/mixin-sdk-go"
)

type (
	Config struct {
		Assets    []*core.Asset  `yaml:"assets"`
		Dapp      mixin.Keystore `json:"dapp"`
		Signers   []*core.Signer `json:"signers"`
		Threshold uint8          `json:"threshold,omitempty"`
	}
)
