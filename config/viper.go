package config

import (
	"github.com/spf13/viper"
)

// Legal o viper, não conhecia.
// Mas nesse caso, que é tão pequeno, não seria melhor usar um .env? Assim não precisa de uma biblioteca só para isso
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
