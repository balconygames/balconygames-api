package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/models"
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

func (s serviceRedisSuite) TestSetScore() {
	score1 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: myUserID,
		},
		LeaderboardID: leaderboardID,
		Value:         100,
	}

	score2 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: topUserID,
		},
		LeaderboardID: leaderboardID,
		Value:         102,
	}

	score3 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: otherUserID1,
		},
		LeaderboardID: leaderboardID,
		Value:         101,
	}

	score4 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: otherUserID2,
		},
		LeaderboardID: leaderboardID,
		Value:         92,
	}

	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()
	_ = repo.SetScore(context.Background(), score1)
	_ = repo.SetScore(context.Background(), score2)
	_ = repo.SetScore(context.Background(), score3)
	_ = repo.SetScore(context.Background(), score4)

	key := scoreKey(score1.LeaderboardID)
	result, err := repo.conn.ZRevRank(context.Background(), key, topUserID).Result()
	s.Require().Nil(err)
	s.Require().Equal(result, int64(0))

	key = scoreKey(score2.LeaderboardID)
	result, err = repo.conn.ZRevRank(context.Background(), key, myUserID).Result()
	s.Require().Nil(err)
	s.Require().Equal(result, int64(2))

	scores, err := repo.ListScores(context.Background(), score1.Scope, leaderboardID)
	s.Require().Nil(err)
	s.Require().Len(scores, 4)
	s.Require().Equal(scores[0].UserID, topUserID)
	s.Require().Equal(scores[0].Position, int64(1))
	s.Require().Equal(scores[1].UserID, otherUserID1)
	s.Require().Equal(scores[1].Position, int64(2))
	s.Require().Equal(scores[2].UserID, myUserID)
	s.Require().Equal(scores[2].Position, int64(3))
	s.Require().Equal(scores[3].UserID, otherUserID2)
	s.Require().Equal(scores[3].Position, int64(4))
}

func (s serviceRedisSuite) TestNoScore() {
	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()

	scores, err := repo.ListScores(context.Background(), sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: myUserID,
	}, leaderboardID)
	s.Require().Nil(err)
	s.Require().Len(scores, 0)
}

func (s serviceRedisSuite) TestMyScore() {
	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	score1 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: myUserID,
		},
		LeaderboardID: leaderboardID,
		Value:         100,
	}

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()
	repo.SetScore(context.Background(), score1)

	scores, err := repo.ListScores(context.Background(), sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: myUserID,
	}, leaderboardID)
	s.Require().Nil(err)
	s.Require().Len(scores, 1)
	s.Require().Equal(scores[0].UserID, myUserID)
}

func (s serviceRedisSuite) TestMyAndTopScore() {
	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	score1 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: myUserID,
		},
		LeaderboardID: leaderboardID,
		Value:         100,
		IP:            "8.8.8.8",
		Country:       "US",
		Timestamp:     1,
	}

	score2 := &models.Score{
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: topUserID,
		},
		LeaderboardID: leaderboardID,
		Value:         101,
		IP:            "1.1.1.1",
		Country:       "BY",
		Timestamp:     2,
	}

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()

	_ = repo.SetScore(context.Background(), score1)
	_ = repo.SetScore(context.Background(), score2)

	scores, err := repo.ListScores(context.Background(), sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: myUserID,
	}, leaderboardID)
	s.Require().Nil(err)
	s.Require().Len(scores, 2)
	s.Require().Equal(scores[0].UserID, topUserID)
	s.Require().Equal(scores[0].IP, "1.1.1.1")
	s.Require().Equal(scores[0].Country, "BY")
	s.Require().Equal(scores[0].Timestamp, int64(2))
	s.Require().Equal(scores[1].UserID, myUserID)
	s.Require().Equal(scores[1].IP, "8.8.8.8")
	s.Require().Equal(scores[1].Country, "US")
	s.Require().Equal(scores[1].Timestamp, int64(1))
}

func (s serviceRedisSuite) TestHugeAmountOfScores() {
	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()

	for i := 0; i < 20; i++ {
		scope := sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: fmt.Sprintf("%d", i*100),
		}

		if i == 5 {
			scope.UserID = myUserID
		}

		score := &models.Score{
			Scope:         scope,
			LeaderboardID: leaderboardID,
			Value:         100.0 + float64(i),
			IP:            "8.8.8.8",
			Country:       "US",
			Timestamp:     1,
		}

		repo.SetScore(context.Background(), score)
	}

	scope := sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: myUserID,
	}

	scores, err := repo.ListScores(context.Background(), scope, leaderboardID)
	s.Require().Nil(err)
	for _, score := range scores {
		if score.UserID == myUserID {
			s.Require().Equal(score.Position, int64(15))
			break
		}
	}
}

func (s serviceRedisSuite) TestTopHugeAmountOfScores() {
	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()

	for i := 0; i < 20; i++ {
		scope := sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: fmt.Sprintf("%d", i*100),
		}

		if i == 5 {
			scope.UserID = myUserID
		}

		score := &models.Score{
			Scope:         scope,
			LeaderboardID: leaderboardID,
			Value:         100.0 + float64(i),
			IP:            "8.8.8.8",
			Country:       "US",
			Timestamp:     1,
		}

		repo.SetScore(context.Background(), score)
	}

	lastUserID := "0"
	topUserID := "1900"

	scope := sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: lastUserID,
	}
	scores, err := repo.ListScores(context.Background(), scope, leaderboardID)
	s.Require().Equal(len(scores), 20)
	s.Require().Nil(err)
	for _, score := range scores {
		if score.UserID == lastUserID {
			s.Require().Equal(score.Position, int64(20))
			break
		}
	}

	scope = sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: topUserID,
	}
	scores, err = repo.ListScores(context.Background(), scope, leaderboardID)
	s.Require().Equal(len(scores), 20)
	s.Require().Nil(err)
	for _, score := range scores {
		if score.UserID == topUserID {
			s.Require().Equal(score.Position, int64(1))
			break
		}
	}
}

func (s serviceRedisSuite) TestMeHugeAmountOfScores() {
	l, err := logging.ConfigForEnv("test").Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	s.Require().Nil(err)

	defer l.Sync()
	logger := l.Sugar()

	repo := NewRedisRepository(s.Conn, logger)
	repo.flush()

	myUserID := "999999"
	topUserID := "000000"
	lastUserID := "111111"

	for i := 0; i < 40; i++ {
		scope := sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
			UserID: fmt.Sprintf("%d", i*100),
		}

		if i == 0 {
			scope.UserID = lastUserID
		}

		if i == 5 {
			scope.UserID = myUserID
		}

		if i == 39 {
			scope.UserID = topUserID
		}

		score := &models.Score{
			Scope:         scope,
			LeaderboardID: leaderboardID,
			Value:         100.0 + float64(i),
			IP:            "8.8.8.8",
			Country:       "US",
			Timestamp:     1,
		}

		repo.SetScore(context.Background(), score)
	}

	scope := sharedmodels.Scope{
		GameID: gameID,
		AppID:  appID,
		UserID: myUserID,
	}
	scores, err := repo.ListScores(context.Background(), scope, leaderboardID)
	s.Require().Equal(len(scores), 40)
	s.Require().Nil(err)

	foundLastUserID := false
	foundTopUserID := false
	foundMyUserID := false
	for _, score := range scores {
		if score.UserID == topUserID {
			s.Require().Equal(score.Position, int64(1))
			foundTopUserID = true
		}

		if score.UserID == myUserID {
			s.Require().Equal(score.Position, int64(35))
			foundMyUserID = true
		}

		if score.UserID == lastUserID {
			s.Require().Equal(score.Position, int64(40))
			foundTopUserID = true
		}
	}

	s.Require().True(foundMyUserID)
	s.Require().True(foundTopUserID)
	s.Require().False(foundLastUserID)
}

// func (s serviceRedisSuite) TestListScoresPRODUCTIONExample() {
// 	l, err := logging.ConfigForEnv("test").Build(
// 		zap.Fields(zap.String("project", "modules/leaderboard")),
// 	)
// 	s.Require().Nil(err)

// 	defer l.Sync()
// 	logger := l.Sugar()

// 	gameID := "d5b0c4c4-483b-425f-9dca-ba1c50bc734f"
// 	appID := "2c7c2c96-b49e-4f29-965f-6b6945cda0bb"
// 	leaderboardID := "bc4378d0-ae8a-4ca7-acfd-db894b406b13"
// 	// oivoodoo
// 	userID := "78ad2eda-bb6c-d404-b522-e14049dfae18"

// 	conn := redis.NewClient(&redis.Options{
// 		Addr: "18.234.247.63:6340",
// 	})
// 	repo := NewRedisRepository(conn, logger)

// 	scope := sharedmodels.Scope{
// 		GameID: gameID,
// 		AppID:  appID,
// 		UserID: userID,
// 	}
// 	scores, err := repo.ListScores(context.Background(), scope, leaderboardID)
// 	s.Require().Equal(len(scores), 20)
// 	s.Require().Nil(err)
// }
