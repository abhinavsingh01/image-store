package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	Secret          string `mapstructure:"SECRET"`
	UserServiceUrl  string `mapstructure:"USER_SVC_URL"`
	AlbumServiceUrl string `mapstructure:"ALBUM_SVC_URL"`
	AuthServiceUrl  string `mapstructure:"AUTH_SVC_URL"`
	ImageServiceUrl string `mapstructure:"IMAGE_SVC_URL"`
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
