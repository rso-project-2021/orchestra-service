package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	StationsAddress string `mapstructure:"stations_address"`
	GrpcAddress     string `mapstructure:"grpc_address"`
	ServerAddress   string `mapstructure:"server_address"`
	GinMode         string `mapstructure:"gin_mode"`
}

// Reads configuration from file or environment variables.
func New(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
