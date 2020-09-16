package handlers

import "gitlab.com/balconygames/analytics/modules/messages/internal/service"

// Handler should define routes for module
type Handler struct {
	service *service.Service
}

// New wrap logic on routing
func New(s *service.Service) *Handler {
	return &Handler{
		service: s,
	}
}
