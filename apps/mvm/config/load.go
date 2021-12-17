package config

import (
	"github.com/fox-one/pkg/config"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

func Load(cfgFile string, cfg *Config) error {
	config.AutomaticLoadEnv("ORACLEEXAMPLE")
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}
	defaultGas(cfg)

	return nil
}

func defaultGas(cfg *Config) {
	if cfg.Gas.Asset.IsNil() {
		cfg.Gas.Asset = uuid.Must(uuid.FromString("965e5c6e-434c-3fa9-b780-c50f43cd955c"))
	}

	if cfg.Gas.Amount.IsZero() {
		cfg.Gas.Amount = decimal.New(1, -8)
	}
}
