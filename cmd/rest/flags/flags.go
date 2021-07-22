package flags

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/symaster1995/ms-starter/pkg/database"
	"time"
)

type ApiConfig struct {
	HttpBindAddress       string
	HttpReadHeaderTimeout time.Duration
	HttpReadTimeout       time.Duration
	HttpWriteTimeout      time.Duration
	Viper                 *viper.Viper
	Domain                string
}

type Config struct {
	Viper    *viper.Viper
	ApiOpts  *ApiConfig
	DBConfig *database.Config
}

func NewConfig(v *viper.Viper) *Config {

	db, err := database.NewConfig(v)

	if err != nil {
		fmt.Println("error setting database config", err)
	}

	api := &ApiConfig{
		HttpBindAddress:       ":8080",
		HttpReadHeaderTimeout: 10 * time.Second,
		HttpReadTimeout:       1 * time.Second,
		HttpWriteTimeout:      1 * time.Second,
		Domain:                "",
	}

	config := &Config{
		Viper: v,
		ApiOpts: api,
		DBConfig: db,
	}

	return config
}