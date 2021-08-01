package config

import (
	"github.com/spf13/viper"
)

type RootConfig struct {
	Server    Server
	Endpoints Endpoints
	Log       Log
}

type Server struct {
	Port int
}

type Log struct {
	Level  string
	Format string
}

type Endpoints struct {
	ServerEndpoint string
}

func NewConfig() *RootConfig {
	viper.SetConfigType("yml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	conf := RootConfig{}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
