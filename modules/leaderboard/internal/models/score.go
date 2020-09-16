package models

import (
	"strconv"

	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// Score should be reported for specific
// leaderboard
type Score struct {
	// add game, app, user ids
	sharedmodels.Scope

	// LeaderboardID should be mapped
	// to specific leaderboard
	LeaderboardID string `json:"leaderboard_id"`

	// User could set name in the game otherwise
	// names could come from social networks.
	Name string `json:"name"`

	// when the user posted score
	Timestamp int64 `json:"timestamp"`

	Value float64 `json:"value"`

	// Type:
	// - me
	// - top
	// - other
	Type string `json:"type"`

	Position int64 `json:"position"`

	// Countries related information to show flag in leaderboard
	IP      string `json:"ip"`
	Country string `json:"country"`
}

func NewScoreByAttrs(scope sharedmodels.Scope, leaderboardID string, attrs map[string]string) (*Score, error) {
	value, err := strconv.ParseFloat(attrs["value"], 64)
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(attrs["timestamp"], 10, 64)
	if err != nil {
		return nil, err
	}

	var position int64 = 0
	if attrs["position"] != "" {
		position, err = strconv.ParseInt(attrs["position"], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	score := &Score{
		Scope: sharedmodels.Scope{
			GameID: scope.GameID,
			AppID:  scope.AppID,
			UserID: attrs["user_id"],
		},
		LeaderboardID: leaderboardID,
		Value:         value,
		Country:       attrs["country"],
		IP:            attrs["ip"],
		Name:          attrs["name"],
		Timestamp:     timestamp,
		Type:          attrs["type"],
		Position:      position,
	}
	return score, nil
}
