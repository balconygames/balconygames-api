package service

import (
	"context"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/models"
)

func (s *Service) CreateLeaderboard(ctx context.Context, leaderboard *models.Leaderboard) error {
	// s.pgRepo.LeaderboardUpsert(ctx, leaderboard)
	return nil
}
