package config

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	DBpath  string `env:"SHORTY_DB_PATH" envDefault:"shorty.db"`
	Address string `env:"SHORTY_ADDRESS" envDefault:":1323"`
}

var (
	config Config
	once   = &sync.Once{}
)

func C() Config {
	once.Do(func() {
		if err := env.Parse(&config); err != nil {
			log.Fatalf("%+v\n", err)
		}
	})

	return config
}
