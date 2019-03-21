package main

import (
	"github.com/spf13/viper"
	"log"
)

var conf *viper.Viper

func initConf() {
	conf = viper.New()
	conf.SetConfigType("yaml")
	conf.SetConfigName("config")
	conf.AddConfigPath("./conf/")
	err := conf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
