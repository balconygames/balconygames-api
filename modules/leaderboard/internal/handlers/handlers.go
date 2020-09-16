package handlers

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/models"
	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/service"
	"gitlab.com/balconygames/analytics/pkg/auth"
	httpreq "gitlab.com/balconygames/analytics/pkg/http"
)

type Handler struct {
	service *service.Service

	logger *zap.SugaredLogger
}

func New(s *service.Service, l *zap.SugaredLogger) *Handler {
	return &Handler{
		service: s,
		logger:  l.With("scope", "leaderboard.handler"),
	}
}

func (h *Handler) CreateLeaderboards(w http.ResponseWriter, r *http.Request) {
	var err error

	data := &models.Leaderboard{}

	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't create leaderboard request"))
		return
	}

	h.service.CreateLeaderboard(r.Context(), data)

	httpreq.NotImplemented(w)
}

func (h *Handler) ListLeaderboards(w http.ResponseWriter, r *http.Request) {
	// TODO: add
	httpreq.NotImplemented(w)
}

type createScoresRequest struct {
	Scores []*models.Score `json:"scores"`
}

// CreateScores is using batch of scores with leaderboard_id per score group
func (h *Handler) CreateScores(w http.ResponseWriter, r *http.Request) {
	var err error

	scope := auth.GetScope(r)

	// block spammer
	if scope.UserID == "a52f61fa-7f00-1d74-6989-474127a981fb" {
		httpreq.NotImplemented(w)
		return
	}

	log := h.logger.With(scope.Fields()...)
	log.Debugf("begin create scores user: %v", scope)

	data := createScoresRequest{}
	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't read leaderboard request"))
		return
	}

	// set ip and country per score
	for _, score := range data.Scores {
		// ensure to have game and app
		// based on JWT token.
		score.UserID = scope.UserID
		score.GameID = scope.GameID
		score.AppID = scope.AppID

		score.IP = r.RemoteAddr
		if score.Country == "" {
			country, err := h.service.Geo.Resolve(score.IP)
			if err != nil {
				httpreq.Error(w, errors.Wrap(err, "can't read geo country by ip address because of internal error"))
				return
			}
			score.Country = country
		}

		log.
			With(score.Scope.Fields()...).
			With("ip", score.IP, "country", score.Country).
			Debugf("refreshed score with geo ip, country information %v", score)
	}

	err = h.service.SetScores(r.Context(), data.Scores)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't upsert leaderboard data"))
		return
	}

	log.Debugf("finished set scores: %v", data.Scores)

	httpreq.OK(w)
}

type listScoresRequest struct {
	LeaderboardIDS []string `json:"leaderboard_id"`
}

type listScoresResponse struct {
	Leaderboards []*models.Leaderboard `json:"leaderboards"`
}

// ListScores respond with list of scores by leaderboard_id's
func (h *Handler) ListScores(w http.ResponseWriter, r *http.Request) {
	var err error

	scope := auth.GetScope(r)

	log := h.logger.With(scope.Fields()...)
	log.Debugf("begin list scores for scope: %v", scope)

	data := listScoresRequest{}
	err = httpreq.Read(r, &data)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't read leaderboard request"))
		return
	}
	log.
		With("leaderboard_id", strings.Join(data.LeaderboardIDS, ",")).
		Infof("list scores for leaderboards: %v", data)

	scores, err := h.service.ListScores(r.Context(), *scope, data.LeaderboardIDS)
	if err != nil {
		httpreq.Error(w, errors.Wrap(err, "can't get scores"))
		return
	}

	httpreq.JSON(w, listScoresResponse{scores})
}
