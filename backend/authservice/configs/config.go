package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port       string `mapstructure:"PORT"`
	UserSvcUrl string `mapstructure:"USER_SVC_URL"`
	Secret     string `mapstructure:"SECRET"`
	LogLevel   string `mapstructure:"LOG_LEVEL"`
}

var AppConfig Config

func LoadConfig() (config *Config) {

	viper.AddConfigPath("./configs/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = viper.Unmarshal(&AppConfig)

	return &AppConfig
}

func GetConfig() (config *Config) {
	return &AppConfig
}
