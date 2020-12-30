package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
	"gitlab.com/balconygames/analytics/modules/auth/internal/service"
	"gitlab.com/balconygames/analytics/pkg/auth"
	httpreq "gitlab.com/balconygames/analytics/pkg/http"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

// Handler have the basic requirements for handlers
type Handler struct {
	service *service.Service

	logger *zap.SugaredLogger
}

// New creates the new instance of Handler
func New(s *service.Service, l *zap.SugaredLogger) *Handler {
	return &Handler{
		service: s,
		logger:  l.With("scope", "auth"),
	}
}

type syncAnomRequest struct {
	// required to be here
	DeviceID string `json:"device_id"`

	// PropertiesSections should be used to get on auth request settings in response
	// for further initialize in the client.
	PropertiesSections []string `json:"props_sections"`
}

type syncAnomResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	JWT      string `json:"jwt"`

	PropertiesSections []*models.Properties `json:"props_sections"`

	// Timestamp could be used to sync the server time
	// and player time in case of using for daily bonuses
	// or other stuff in the game.
	Timestamp int64 `json:"timestamp"`
}

type syncRegRequest struct {
	// required to be here
	DeviceID string `json:"device_id"`
	GuestID  string `json:"guest_id"`

	Name  string `json:"name"`
	Email string `json:"email"`

	Network   string `json:"network"`
	NetworkID string `json:"network_id"`

	// PropertiesSections should be used to get on auth request settings in response
	// for further initialize in the client.
	PropertiesSections []string `json:"props_sections"`
}

type syncRegResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	GuestID  string `json:"guest_id"`

	NetworkID string `json:"network_id"`
	Network   string `json:"network"`

	JWT string `json:"jwt"`

	PropertiesSections []*models.Properties `json:"props_sections"`

	// Timestamp could be used to sync the server time
	// and player time in case of using for daily bonuses
	// or other stuff in the game.
	Timestamp int64 `json:"timestamp"`
}

func timestamp() int64 {
	return time.Now().Unix()
}

// SyncAnomHandler should start collecting the users and build user_archive records
// it should have user_id, device_id, email, network_id
func (h *Handler) SyncAnomHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	gameID := chi.URLParam(r, "game_id")
	appID := chi.URLParam(r, "app_id")

	log := h.logger.With("game_id", gameID, "app_id", appID)

	data := syncAnomRequest{}
	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't read sync body"))
		return
	}

	// // block spammer
	// if data.DeviceID == "e4f5142c7553c1bddefaee1e3cc00d1e" {
	// 	httpreq.NotImplemented(w)
	// 	return
	// }

	// sign up request should converted
	// to models user and synced up
	user := models.User{
		DeviceID: data.DeviceID,
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
		},
	}

	log = log.With("device_id", data.DeviceID)
	log.Debug("begin auth sync request")

	properties, err := h.service.Sync(r.Context(), data.PropertiesSections, &user)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't sync the user"))
		return
	}

	signer := auth.GetSigner(r)

	jwt, err := signer.Encode(auth.Claims{
		UserInfo: auth.UserInfo{
			AppID:    user.Scope.AppID,
			GameID:   user.Scope.GameID,
			UserID:   user.Scope.UserID,
			DeviceID: user.DeviceID,
			Type:     auth.GuestType,
		},
	})
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't generate jwt token properly"))
		return
	}

	log.Debugf("respond with properties: %v", properties)

	out := &syncAnomResponse{
		PropertiesSections: properties,
		UserID:             user.Scope.UserID,
		UserName:           user.Name,
		JWT:                jwt,
		Timestamp:          timestamp(),
	}

	httpreq.JSON(w, out)
}

// SyncRegHandler should start collecting the users and build user_archive records
// it should have user_id, device_id, email, network_id
func (h *Handler) SyncRegHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	gameID := chi.URLParam(r, "game_id")
	appID := chi.URLParam(r, "app_id")

	log := h.logger.With("game_id", gameID, "app_id", appID)

	data := syncRegRequest{}
	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't read sync body"))
		return
	}

	// block spammer
	// if data.DeviceID == "e4f5142c7553c1bddefaee1e3cc00d1e" {
	// 	httpreq.NotImplemented(w)
	// 	return
	// }

	// sign up request should converted
	// to models user and synced up
	user := models.User{
		DeviceID:  data.DeviceID,
		Email:     data.Email,
		Name:      data.Name,
		Network:   data.Network,
		NetworkID: data.NetworkID,
		GuestID:   data.GuestID,
		Scope: sharedmodels.Scope{
			GameID: gameID,
			AppID:  appID,
		},
	}

	log = log.With("device_id", data.DeviceID)
	log.Debug("begin auth sync request")

	properties, err := h.service.UsersSync(r.Context(), data.PropertiesSections, &user)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't sync the user"))
		return
	}

	signer := auth.GetSigner(r)

	jwt, err := signer.Encode(auth.Claims{
		UserInfo: auth.UserInfo{
			AppID:    user.Scope.AppID,
			GameID:   user.Scope.GameID,
			UserID:   user.Scope.UserID,
			DeviceID: user.DeviceID,
			Type:     auth.RealType,
		},
	})
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't generate jwt token properly"))
		return
	}

	log.Debugf("respond with properties: %v", properties)

	out := &syncRegResponse{
		PropertiesSections: properties,

		UserID:   user.Scope.UserID,
		UserName: user.Name,
		GuestID:  user.GuestID,

		JWT: jwt,

		Network:   user.Network,
		NetworkID: user.NetworkID,

		Timestamp: timestamp(),
	}

	httpreq.JSON(w, out)
}
