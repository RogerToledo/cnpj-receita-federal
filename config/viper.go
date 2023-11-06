package config

import (
	"github.com/spf13/viper"
)

func Load() (*viper.Viper, error) {
	conf := viper.GetViper()
	conf.AddConfigPath("./config")
	conf.SetConfigName("conf")
	conf.SetConfigType("yaml")
	if err := conf.ReadInConfig(); err != nil {
		return nil, err
	}

	return conf, nil
}
