package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	ConnUrl string
}

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("yml")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
