package models

import (
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// Pixel passed in requests and then transmitted to message queue
type Pixel struct {
	sharedmodels.Scope

	Data map[string]string `json:"data"`
}
