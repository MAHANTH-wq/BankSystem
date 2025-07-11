package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource      string `mapstructure:"DBSOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // we can use other formats like "json", "yaml", etc.
	viper.AutomaticEnv()       // read environment variables that match the keys

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return
}
