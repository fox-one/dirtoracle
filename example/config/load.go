package config

import (
	"github.com/fox-one/pkg/config"
)

func Load(cfgFile string, cfg *Config) error {
	config.AutomaticLoadEnv("ORACLEEXAMPLE")
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}

	return nil
}
