package models

import (
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// Leaderboard is container for scores list
// Example: Coins, Levels highscores
type Leaderboard struct {
	sharedmodels.Scope

	ID   string `json:"id"`
	Name string `json:"name"`

	Scores []*Score `json:"scores"`
}
