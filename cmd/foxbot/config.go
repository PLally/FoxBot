package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupConfig() {
	viper.SetEnvPrefix("FOX_BOT")
	viper.SetConfigName("foxbot_config")

	viper.AutomaticEnv()

	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/foxbot/")
	viper.AddConfigPath("$HOME/.foxbot")
	viper.AddConfigPath(".")

	viper.SetDefault("prefix", ">")
	viper.SetDefault("logfile", "foxbot.log")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
