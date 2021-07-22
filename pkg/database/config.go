package database

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DBName            string
	User              string
	Host              string
	Port              string
	Password          string
	ConnectionTimeout int
	URL               string
}

func NewConfig(v *viper.Viper) (*Config, error) {

	config, err := setConfig(v)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func setConfig(v *viper.Viper) (*Config, error) {
	v.SetConfigName("db_config")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetConfigType("yml")

	var config Config
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetDefault("dbname", "ms_starter_demo")
	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	config.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.User, config.Password, config.Host, config.Port, config.DBName)

	return &config, nil
}
