package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
	"gitlab.com/balconygames/analytics/pkg/logging"
	test_helpers "gitlab.com/balconygames/analytics/pkg/test_helpers"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type serviceRedisSuite struct {
	test_helpers.RedisSuite

	logger *zap.SugaredLogger
}

func TestRedisEntrypointSuite(t *testing.T) {
	handler := &serviceRedisSuite{
		RedisSuite: test_helpers.NewDefaultRedisSuite(t, "test_leaderboard"),
		logger:     zaptest.NewLogger(t).Sugar(),
	}

	suite.Run(t, handler)
}

func (s serviceRedisSuite) TestGetProperties() {
	var err error

	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	repo := NewRedisRepository(s.Conn, logger)

	scope := sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: userID,
	}
	collection := []*models.Properties{
		{
			Section: "users",
			Data:    map[string]string{"users1": "users-value1", "users2": "users-value2"},
			Scope:   scope,
		},
		{
			Section: "game",
			Data:    map[string]string{"game1": "game-value1", "game2": "game-value2"},
			Scope:   scope,
		},
	}
	err = repo.SetProperties(context.Background(), collection)
	s.Require().Nil(err)

	output, err := repo.GetProperties(context.Background(), []string{"users", "cloud", "game"}, &scope)
	s.Require().Nil(err)
	s.Require().Len(output, 2)
	s.Require().Equal("users-value1", output[0].Data["users1"])
	s.Require().Equal("users-value2", output[0].Data["users2"])
	s.Require().Equal("game-value1", output[1].Data["game1"])
	s.Require().Equal("game-value2", output[1].Data["game2"])
}
