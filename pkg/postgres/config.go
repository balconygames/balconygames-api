package postgres

import "fmt"

type Config struct {
	User       string `envconfig:"DB_USER"`
	Pass       string `envconfig:"DB_PASS"`
	Host       string `envconfig:"DB_HOST"`
	Name       string `envconfig:"DB_NAME"`
	DisableSSL bool   `envconfig:"DB_DISABLE_SSL"`
}

func (c Config) URL() string {
	var sslmode string
	if c.DisableSSL {
		sslmode = "sslmode=disable"
	}
	return fmt.Sprintf("postgres://%s:%s@%s/%s?%s", c.User, c.Pass, c.Host, c.Name, sslmode)
}
