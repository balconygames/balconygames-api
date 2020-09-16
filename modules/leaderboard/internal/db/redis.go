package db

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/models"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type RedisRepository struct {
	conn *redis.Client

	logger *zap.SugaredLogger
}

func NewRedisRepository(c *redis.Client, l *zap.SugaredLogger) *RedisRepository {
	return &RedisRepository{c, l.With("scope", "redis")}
}

func scoreKey(leaderboardID string) string {
	return fmt.Sprintf("scores:%s", leaderboardID)
}

func userKey(leaderboardID, userID string) string {
	return fmt.Sprintf("users:%s:%s", leaderboardID, userID)
}

// internal function to clean up the database before to run
// test cases.
func (r *RedisRepository) flush() error {
	return r.conn.FlushAll(context.Background()).Err()
}

// SetScore should set score by scope.
func (r *RedisRepository) SetScore(ctx context.Context, score *models.Score) error {
	var err error

	pipe := r.conn.Pipeline()

	_, err = pipe.ZAdd(
		ctx,
		scoreKey(score.LeaderboardID),
		&redis.Z{Score: score.Value, Member: score.UserID},
	).Result()
	if err != nil {
		return err
	}

	var values []string
	values = append(values, "country", score.Country)
	values = append(values, "ip", score.IP)
	values = append(values, "name", score.Name)
	values = append(values, "timestamp", strconv.FormatInt(score.Timestamp, 10))
	values = append(values, "value", strconv.FormatFloat(score.Value, 'f', -1, 64))
	values = append(values, "user_id", score.UserID)

	r.logger.Debugf("set score attributes: %v", values)

	_, err = pipe.HSet(ctx, userKey(score.LeaderboardID, score.UserID), values).Result()
	if err != nil {
		return err
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

const minRange = 100

func (r *RedisRepository) ListScores(ctx context.Context, scope sharedmodels.Scope, leaderboardID string) ([]*models.Score, error) {
	var err error

	r.logger.Debug("get current user rank ", "user id: ", scope.UserID, " leaderboard id: ", leaderboardID)

	myRank, _ := r.conn.ZRevRank(ctx, scoreKey(leaderboardID), scope.UserID).Result()

	r.logger.Debug("got current user rank ", "rank ", myRank)

	r.logger.Debug("begin pipeline to get other ranked users and top user by leaderboard id: ", leaderboardID, " user id: ", scope.UserID)
	// begin get ids
	pipe := r.conn.Pipeline()
	pipe.ZRevRange(ctx, scoreKey(leaderboardID), 0, 0)

	beginMarkerRank := myRank - 2*minRange
	if beginMarkerRank < 0 {
		beginMarkerRank = 0
	}
	pipe.ZRevRange(ctx, scoreKey(leaderboardID), beginMarkerRank, myRank+2*minRange)
	results, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error(errors.WithMessage(err, "can't run pipeline to receive other users and top ranked user"))
		return nil, err
	}

	topSlice, _ := results[0].(*redis.StringSliceCmd)
	top, err := topSlice.Result()

	usersSlice, _ := results[1].(*redis.StringSliceCmd)
	users, err := usersSlice.Result()

	r.logger.Debug("got top user: ", top, " and users: ", len(users))

	// end get ids

	// begin get attributes
	pipe = r.conn.Pipeline()
	pipe.HGetAll(ctx, userKey(leaderboardID, scope.UserID))
	if len(top) > 0 {
		pipe.HGetAll(ctx, userKey(leaderboardID, top[0]))
	}
	for _, user := range users {
		pipe.HGetAll(ctx, userKey(leaderboardID, user))
	}
	results, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var scores []*models.Score

	for _, res := range results {
		if cmd, ok := res.(*redis.StringStringMapCmd); ok {
			attrs, err := cmd.Result()
			if err != nil {
				err = errors.WithMessage(err, "can't get other user score")
				r.logger.Error(err)
				return nil, err
			}
			if len(attrs) == 0 {
				continue
			}

			r.logger.Debugf("get score attributes: %v", attrs)

			score, err := models.NewScoreByAttrs(scope, leaderboardID, attrs)
			if err != nil {
				err = errors.WithMessage(err, "can't build other user score")
				r.logger.Error(err)
				return nil, err
			}
			// my
			if len(scores) > 0 && score.UserID == scores[0].UserID {
				continue
			}
			//top
			if len(scores) > 1 && score.UserID == scores[1].UserID {
				continue
			}

			// should skip users my, top from other lists
			scores = append(scores, score)
		}
	}

	// add positions
	pipe = r.conn.Pipeline()

	var scorePositions []*models.Score

	for _, score := range scores {
		pipe.ZRevRank(ctx, scoreKey(leaderboardID), score.UserID)
		scorePositions = append(scorePositions, score)
	}

	results, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for i, res := range results {
		if cmd, ok := res.(*redis.IntCmd); ok {
			rank, err := cmd.Result()
			if err != nil {
				return nil, err
			}
			scorePositions[i].Position = rank + 1
		}
	}

	// end get attributes
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Position < scores[j].Position
	})

	var myScoreIndex int
	for i, score := range scores {
		if score.UserID == scope.UserID {
			myScoreIndex = i
			break
		}
	}

	var newScores []*models.Score
	// 38 > 20
	if myScoreIndex > minRange {
		newScores = append(newScores, scores[myScoreIndex])
		for i := 1; len(newScores) < minRange; i++ {
			if myScoreIndex+i < len(scores) {
				newScores = append(newScores, scores[myScoreIndex+i])
			}
			if myScoreIndex-i > 0 {
				newScores = append(newScores, scores[myScoreIndex-i])
			}
		}

		newScores = append(newScores, scores[0])
		scores = newScores

		// end get attributes
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Position < scores[j].Position
		})
	}

	if len(scores) >= minRange {
		scores = scores[0:minRange]
	}

	for _, score := range scores {
		if score.Position == 1 {
			score.Type = "top"
		} else if score.UserID == scope.UserID {
			score.Type = "me"
		} else {
			score.Type = "other"
		}
	}

	return scores, nil
}
