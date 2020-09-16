package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	User string `envconfig:"USER"`
	Pass string `envconfig:"PASS"`
	Host string `envconfig:"HOST"`
	Port int    `envconfig:"PORT"`
	Name string `envconfig:"NAME"`
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Config) Options() *redis.Options {
	options := &redis.Options{
		Addr: c.Addr(),
	}

	if c.User != "" {
		options.Username = c.User
	}

	if c.Pass != "" {
		options.Password = c.Pass
	}

	return options
}
