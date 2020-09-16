package service

import (
	"context"

	"go.uber.org/zap"

	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type PostgresRepository interface {
	GetAppInfo(ctx context.Context, app *sharedmodels.App) error
}

// Service contains all dependencies to perform common service tasks.
type Service struct {
	pgRepo PostgresRepository

	logger *zap.SugaredLogger
}

// NewService should build the service to having the layer between handlers
// and repositories.
func NewService(r PostgresRepository, l *zap.SugaredLogger) *Service {
	return &Service{
		pgRepo: r,
		logger: l,
	}
}

func (s *Service) GetAppInfo(ctx context.Context, clientApp *sharedmodels.App) error {
	return s.pgRepo.GetAppInfo(ctx, clientApp)
}
