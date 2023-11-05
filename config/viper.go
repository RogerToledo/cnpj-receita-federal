package config

import (
	"github.com/spf13/viper"
)

func Load() (*viper.Viper, error) {
	conf := viper.GetViper()
	conf.AddConfigPath(".")
	conf.SetConfigFile("config")
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	return conf, nil
}
