package models

import (
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// Network could be FACEBOOK, GOOGLE, APPLE
// once we matched user profiles, we should start
// to use JWT as signed user.
type Network struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// User should be synced with database
// and store inside of `users` table.
// TODO: should we rename to Guest?
type User struct {
	sharedmodels.Scope

	Name    string `json:"user_name"`
	Email   string `json:"email"`
	GuestID string `json:"guest_id"`

	// required to be here
	DeviceID string `json:"device_id"`

	// TODO: Networks should be associated with user, could have multiple
	// network profiles per user.
	// Networks []Network

	Network   string `json:"network"`
	NetworkID string `json:"network_id"`
}

// Properties should store settings per game and app
type Properties struct {
	sharedmodels.Scope

	Section string `json:"section"`

	Data map[string]string `json:"data"`
}
