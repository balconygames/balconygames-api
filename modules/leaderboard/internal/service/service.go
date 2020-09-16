package service

import (
	"context"

	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/models"
	"gitlab.com/balconygames/analytics/pkg/geo"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type PostgresRepository interface {
}

type RedisRepository interface {
	SetScore(context.Context, *models.Score) error
	ListScores(context.Context, sharedmodels.Scope, string) ([]*models.Score, error)
}

// Service contains all dependencies to perform common service tasks.
type Service struct {
	pgRepo    PostgresRepository
	redisRepo RedisRepository

	Geo *geo.DB

	logger *zap.SugaredLogger
}

// NewService should build the service to having the layer between handlers
// and repositories.
func NewService(r PostgresRepository, rp RedisRepository, g *geo.DB, l *zap.SugaredLogger) *Service {
	return &Service{
		pgRepo:    r,
		redisRepo: rp,
		logger:    l,
		Geo:       g,
	}
}
