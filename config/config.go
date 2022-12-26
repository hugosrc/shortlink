package config

import (
	"log"

	"github.com/spf13/viper"
)

func Init() viper.Viper {
	config := viper.New()

	config.AddConfigPath(".")
	config.SetConfigFile(".env")
	if err := config.ReadInConfig(); err != nil {
		log.Println(err)
	}

	config.AutomaticEnv()

	return *config
}
