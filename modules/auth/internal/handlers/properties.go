package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
	"gitlab.com/balconygames/analytics/pkg/auth"
	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

type setPropertiesRequest struct {
	PropertiesSections []*models.Properties `json:"props_sections"`
}

// SetPropertiesHandler set properties list with their section names for
// user by JWT token.
func (h *Handler) SetPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	scope := auth.GetScope(r)

	data := setPropertiesRequest{}
	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't read sync body"))
		return
	}

	log := h.logger.With(scope.Fields()...)
	// ensure to have scope by JWT user
	for _, properties := range data.PropertiesSections {
		properties.AppID = scope.AppID
		properties.GameID = scope.GameID
		properties.UserID = scope.UserID
	}
	log.Debugf("begin set properties %v", data.PropertiesSections)

	err = h.service.SetProperties(r.Context(), data.PropertiesSections)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't set properties for the user"))
		return
	}

	log.Debugf("done set properties %v", data.PropertiesSections)

	httpreq.OK(w)
}

type getPropertiesRequest struct {
	PropertiesSections []string `json:"props_sections"`
}

type getPropertiesResponse struct {
	PropertiesSections []*models.Properties `json:"props_sections"`
}

// GetPropertiesHandler respond with properties list by passed section name
// for user by JWT token.
func (h *Handler) GetPropertiesHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	scope := auth.GetScope(r)

	log := h.logger.With(scope.Fields()...)
	log.Debug("begin get properties request")

	data := getPropertiesRequest{}
	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't read sync body"))
		return
	}

	log.Debugf("done reading request data: %v", data)

	properties, err := h.service.GetProperties(r.Context(), data.PropertiesSections, scope)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't get properties the user"))
		return
	}

	log.Debugf("received properties list: %v", properties)

	httpreq.JSON(w, getPropertiesResponse{properties})
}
