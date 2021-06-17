package flags

import (
	"github.com/spf13/viper"
	"time"
)

type ApiOpts struct {
	HttpBindAddress       string
	HttpReadHeaderTimeout time.Duration
	HttpReadTimeout       time.Duration
	HttpWriteTimeout      time.Duration
	Viper                 *viper.Viper
	Domain                string
}

func NewOpts(v *viper.Viper) *ApiOpts {
	return &ApiOpts{
		Viper:                 v,
		HttpBindAddress:       ":8080",
		HttpReadHeaderTimeout: 10 * time.Second,
		HttpReadTimeout:       1 * time.Second,
		HttpWriteTimeout:      1 * time.Second,
		Domain:                "",
	}
}
