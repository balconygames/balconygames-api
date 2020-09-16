package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/primary/internal/service"
	httpreq "gitlab.com/balconygames/analytics/pkg/http"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type Handler struct {
	service *service.Service

	logger *zap.SugaredLogger
}

func New(s *service.Service, l *zap.SugaredLogger) *Handler {
	return &Handler{
		service: s,
		logger:  l.With("scope", "primary.handler"),
	}
}

func (h *Handler) GetAppInfo(w http.ResponseWriter, r *http.Request) {
	var err error

	app := &sharedmodels.App{}
	app.AppID = chi.URLParam(r, "app_id")
	app.GameID = chi.URLParam(r, "game_id")

	err = h.service.GetAppInfo(r.Context(), app)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't process get app info request using service"))
		return
	}

	// respond with modified app
	httpreq.JSON(w, app)
}
