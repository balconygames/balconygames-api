package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// TODO: move to separate service via grpc calls to fetch properties per user
type RedisRepository struct {
	conn *redis.Client

	logger *zap.SugaredLogger
}

func NewRedisRepository(c *redis.Client, l *zap.SugaredLogger) *RedisRepository {
	return &RedisRepository{c, l.With("scope", "redis")}
}

func propsKey(section string, scope sharedmodels.Scope) string {
	return fmt.Sprintf("props:%s:%s", section, scope.UserID)
}

// GetProperties should respond with user properties per game, app, user id and section
// Example: cloud properties, game properties and so on
func (r *RedisRepository) GetProperties(ctx context.Context, sections []string, scope *sharedmodels.Scope) ([]*models.Properties, error) {
	log := r.logger.With(scope.Fields()...)

	var result []*models.Properties

	pipe := r.conn.Pipeline()

	for _, section := range sections {
		key := propsKey(section, *scope)
		log.
			With("key", key).
			With("section", section).
			Debug("run pipe to get properties")
		pipe.HGetAll(ctx, key).Result()
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return []*models.Properties{}, errors.WithMessage(err, "can't exec pipeline to get properties from redis")
	}

	for i, res := range cmds {
		if cmd, ok := res.(*redis.StringStringMapCmd); ok {
			attrs, err := cmd.Result()
			log.
				With("cmd-name", cmd.FullName()).
				With("section", sections[i]).
				Debugf("got properties from redis: %v", attrs)
			if err != nil {
				continue
			}
			if len(attrs) == 0 {
				continue
			}

			properties := &models.Properties{}
			properties.Data = attrs
			properties.Scope = *scope
			properties.Section = sections[i]

			result = append(result, properties)
		}
	}

	return result, nil
}

// SetProperties should respond with user properties per game, app, user id and section
func (r *RedisRepository) SetProperties(ctx context.Context, collection []*models.Properties) error {
	pipe := r.conn.Pipeline()

	for _, properties := range collection {
		var values []string
		for key, value := range properties.Data {
			values = append(values, key, value)
		}

		key := propsKey(properties.Section, properties.Scope)
		r.logger.
			With("key", key).
			With("section", properties.Section).
			Debugf("set values: %v", values)
		pipe.HSet(ctx, key, values)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return errors.WithMessage(err, "can't exec pipeline to set properties from redis")
	}

	return nil
}
