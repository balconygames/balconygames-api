package service

import (
	"context"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/models"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

func (s Service) SetScores(ctx context.Context, scores []*models.Score) error {
	for _, score := range scores {
		err := s.SetScore(ctx, score.Scope, score)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Service) SetScore(ctx context.Context, scope sharedmodels.Scope, score *models.Score) error {
	return s.redisRepo.SetScore(ctx, score)
}

func (s Service) ListScores(ctx context.Context, scope sharedmodels.Scope, leaderboardID []string) ([]*models.Leaderboard, error) {
	var ls []*models.Leaderboard

	for _, id := range leaderboardID {
		scores, err := s.redisRepo.ListScores(ctx, scope, id)
		if err != nil {
			return nil, err
		}

		ls = append(ls, &models.Leaderboard{
			ID:     id,
			Scope:  scope,
			Scores: scores,
		})
	}

	return ls, nil
}
