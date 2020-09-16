package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// PostgresRepository should define accessible methods
// to database layer.
type PostgresRepository interface {
	AnomSync(context.Context, *models.User) error
	UsersSync(context.Context, *models.User) error
}

// RedisRepository should provide access to user properties
type RedisRepository interface {
	GetProperties(ctx context.Context, sections []string, scope *sharedmodels.Scope) ([]*models.Properties, error)
	SetProperties(ctx context.Context, collection []*models.Properties) error
}

// Service contains all dependencies to perform common service tasks.
type Service struct {
	repoPG    PostgresRepository
	repoRedis RedisRepository

	logger *zap.SugaredLogger
}

// NewService should build the service to having the layer between handlers
// and repositories.
func NewService(r PostgresRepository, rp RedisRepository, l *zap.SugaredLogger) *Service {
	return &Service{
		repoPG:    r,
		repoRedis: rp,
		logger:    l,
	}
}

// Sync the user with the database, in case if we have the record
// we should respond with properties for the user to boot in the app.
func (s *Service) Sync(ctx context.Context, propertiesSections []string, user *models.User) ([]*models.Properties, error) {
	err := s.repoPG.AnomSync(ctx, user)
	if err != nil {
		return nil, errors.WithMessage(err, "can't sync anonymous user with device id")
	}

	if len(propertiesSections) == 0 {
		return nil, nil
	}

	properties, err := s.repoRedis.GetProperties(ctx, propertiesSections, &user.Scope)
	if err != nil {
		return nil, errors.WithMessage(err, "can't get properties list")
	}

	return properties, nil
}

// UsersSync should create or update the record in database for users signed in
// via network FACEBOOK, GOOGLE, APPLE.
func (s *Service) UsersSync(ctx context.Context, propertiesSections []string, user *models.User) ([]*models.Properties, error) {
	err := s.repoPG.UsersSync(ctx, user)
	if err != nil {
		return nil, errors.WithMessage(err, "can't sync network user with device id, network and network id")
	}

	if len(propertiesSections) == 0 {
		return nil, nil
	}

	properties, err := s.repoRedis.GetProperties(ctx, propertiesSections, &user.Scope)
	if err != nil {
		return nil, errors.WithMessage(err, "can't get properties list")
	}

	return properties, nil
}

func (s *Service) GetProperties(ctx context.Context, propertiesSections []string, scope *sharedmodels.Scope) ([]*models.Properties, error) {
	properties, err := s.repoRedis.GetProperties(ctx, propertiesSections, scope)
	if err != nil {
		return nil, errors.WithMessage(err, "can't get properties list")
	}
	if properties == nil {
		return []*models.Properties{}, nil
	}

	return properties, nil
}

func (s *Service) SetProperties(ctx context.Context, collection []*models.Properties) error {
	err := s.repoRedis.SetProperties(ctx, collection)
	if err != nil {
		return errors.WithMessage(err, "can't set properties list")
	}

	return nil
}
