package handlers

import (
	"net/http"

	"github.com/pkg/errors"

	"gitlab.com/balconygames/analytics/modules/pixel/internal/models"
	"gitlab.com/balconygames/analytics/modules/pixel/internal/service"
	"gitlab.com/balconygames/analytics/pkg/auth"
	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

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

// PixelHandler should handle request with json body
// Requiremeents:
// - should require to have device-id
// - should require to have user-id
// - should require to have parameters as dictionary
func (h *Handler) PixelHandler(w http.ResponseWriter, r *http.Request) {
	scope := auth.GetScope(r)

	pr := models.Pixel{Scope: *scope}

	err := httpreq.Read(r, &pr)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't process pixel request"))
		return
	}

	err = h.service.Push(pr)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't push message"))
		return
	}

	httpreq.OK(w)
}
