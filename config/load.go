package config

import (
	"github.com/fox-one/pkg/config"
)

func Load(cfgFile string, cfg *Config) error {
	config.AutomaticLoadEnv("DIRTORACLE")
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}

	defaultVote(cfg)
	return nil
}
