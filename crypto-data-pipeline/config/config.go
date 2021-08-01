package config

import "github.com/spf13/viper"

type RootConfig struct {
	Server   Server
	Database Database
	Log      Log
}

type Server struct {
	Http Http
	Grpc Grpc
}

type Http struct {
	Port  string
	DbLog bool
}

type Grpc struct {
	Port string
}

type Database struct {
	Endpoint string
	Port     int
	Username string
	Password string
	DbName   string
}

type Log struct {
	Level  string
	Format string
}

var Config RootConfig

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
