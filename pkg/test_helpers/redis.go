package test_helpers

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"

	redisconf "gitlab.com/balconygames/analytics/pkg/redis"
)

type RedisSuite struct {
	suite.Suite

	Conn *redis.Client

	Config redisconf.Config
}

func NewDefaultRedisSuite(t *testing.T, dbName string) RedisSuite {
	// Default settings based on docker-compose.yml
	suite := RedisSuite{
		Config: redisconf.Config{
			Port: 6340,
			Host: "127.0.0.1",
		},
	}

	suite.Conn = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			suite.Config.Host,
			suite.Config.Port,
		),
	})

	return suite
}

func (s *RedisSuite) drop(ctx context.Context) {
	s.Conn.FlushAll(ctx)
}

func (s *RedisSuite) create(ctx context.Context) {
}

func (s *RedisSuite) SetupTest() {
	ctx := context.Background()

	s.drop(ctx)
	s.create(ctx)
}

func (s *RedisSuite) TearDownTest() {
	ctx := context.Background()

	s.drop(ctx)
}
