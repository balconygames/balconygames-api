package models

// Scope should be inserted to
// the models that's depending
// on game, app, user
type Scope struct {
	GameID string `json:"game_id"`
	AppID  string `json:"app_id"`
	UserID string `json:"user_id"`
}

// Fields returns the list of fields for logger
func (s Scope) Fields() []interface{} {
	return []interface{}{"game_id", s.GameID, "app_id", s.AppID, "user_id", s.UserID}
}
