package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

var c *Config

type Config struct {
	Port             int    `envconfig:"PORT" default:"18080"`
	ProxyDomain      string `envconfig:"PROXY_DOMAIN" default:"localhost"`
	OriginPort       int    `envconfig:"ORIGIN_PORT" default:"443"`
	OriginBaseDomain string `envconfig:"ORIGIN_BASE_DOMAIN"`
	OriginScheme     string `envconfig:"ORIGIN_SCHEME" default:"http"`
}

func LoadConfig() *Config {
	var cnf Config
	if c != nil {
		return c
	}
	err := envconfig.Process("", &cnf)
	if err != nil {
		panic(err)
	}
	c = &cnf
	return c
}

func GetConfig() *Config {
	return c
}
