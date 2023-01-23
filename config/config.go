package config

import (
	"github.com/hugosrc/shortlink/internal/util"
	"github.com/spf13/viper"
)

func Init() (*viper.Viper, error) {
	config := viper.New()

	config.AddConfigPath(".")
	config.SetConfigFile(".env")
	if err := config.ReadInConfig(); err != nil {
		return nil, util.WrapErrorf(err, util.ErrCodeUnknown, "error")
	}

	config.AutomaticEnv()

	return config, nil
}
