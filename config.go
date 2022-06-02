package main

import (
	"github.com/spf13/viper"
)

var (
	conf Config
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./deploy")
	viper.AddConfigPath("/jike-bot/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	WhAPIKey string `mapstructure:"wh_api_key"`
	PeAPIKey string `mapstructure:"pe_api_key"`

	CachePath string `mapstructure:"cache_path"`
	AreaCode  string `mapstructure:"area_code"`
	Phone     string `mapstructure:"phone"`
	Password  string `mapstructure:"password"`
	Spec      string `mapstructure:"spec"`

	UnsplashAK string `mapstructure:"unsplash_ak"`
	UnsplashSK string `mapstructure:"unsplash_sk"`
}
