package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"gitlab.com/balconygames/analytics/modules/auth/internal/db"
	"gitlab.com/balconygames/analytics/modules/auth/internal/service"
	"gitlab.com/balconygames/analytics/pkg/runtime"
	test_helpers "gitlab.com/balconygames/analytics/pkg/test_helpers"
)

type serviceSuite struct {
	test_helpers.PostgresSuite
	redisSuite test_helpers.RedisSuite

	logger *zap.SugaredLogger
}

func TestEntrypointSuite(t *testing.T) {
	handler := &serviceSuite{
		PostgresSuite: test_helpers.NewDefaultPostgresSuite(t, "test_auth"),
		redisSuite:    test_helpers.NewDefaultRedisSuite(t, "test_auth"),
		logger:        zaptest.NewLogger(t).Sugar(),
	}

	suite.Run(t, handler)
}

type signupAnomRequestTest struct {
	DeviceID string `json:"device_id"`
}

type signupAnomResponseTest struct {
	UserID string `json:"user_id"`

	JWT string `json:"jwt"`
}

type signupRealRequestTest struct {
	DeviceID string `json:"device_id"`
	GuestID  string `json:"guest_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`

	Network   string `json:"network"`
	NetworkID string `json:"network_id"`
}

type signupRealResponseTest struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	GuestID  string `json:"guest_id"`

	JWT string `json:"jwt"`

	Network   string `json:"network"`
	NetworkID string `json:"network_id"`
}

func (s *serviceSuite) TestRequestHandlerToVerifyToken() {
	var err error

	redisRepo := db.NewRedisRepository(s.redisSuite.Conn, s.logger)
	repo := db.NewPostgresRepository(s.PostgresPool)

	svc := service.NewService(repo, redisRepo, s.logger)

	h := New(svc, s.logger)

	spec := runtime.Spec{
		Env:             "test",
		JWTClientSecret: "<secret>",
	}
	r := runtime.New("api", spec)
	s.Require().NoError(err)

	var router chi.Router

	r.WithRoutes(func(r1 chi.Router) {
		router = r1

		r.WithClientTokenSigner(r1, func(r2 chi.Router) {
			r2.Post("/auth/v1/games/{game_id}/apps/{app_id}/anonymous/sync", h.SyncAnomHandler)
		})
	})

	reqBody, err := json.Marshal(&signupAnomRequestTest{
		DeviceID: "0000-0000-0000",
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/v1/games/1/apps/2/anonymous/sync", bytes.NewBuffer(reqBody))
	router.ServeHTTP(w, req)
	s.Require().Equal(200, w.Code)

	respBody := w.Body.Bytes()
	resp := &signupAnomResponseTest{}
	err = json.Unmarshal(respBody, &resp)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.UserID)
}

func (s *serviceSuite) TestUserWorkflow() {
	var err error

	redisRepo := db.NewRedisRepository(s.redisSuite.Conn, s.logger)
	repo := db.NewPostgresRepository(s.PostgresPool)

	svc := service.NewService(repo, redisRepo, s.logger)

	h := New(svc, s.logger)

	spec := runtime.Spec{
		Env:             "test",
		JWTClientSecret: "<secret>",
	}
	r := runtime.New("api", spec)
	s.Require().NoError(err)

	var router chi.Router

	r.WithRoutes(func(r1 chi.Router) {
		router = r1

		r.WithClientTokenSigner(r1, func(r2 chi.Router) {
			r2.Post("/auth/v1/games/{game_id}/apps/{app_id}/anonymous/sync", h.SyncAnomHandler)
			r2.Post("/auth/v1/games/{game_id}/apps/{app_id}/users/sync", h.SyncRegHandler)
		})
	})

	reqAnomBody, err := json.Marshal(&signupAnomRequestTest{
		DeviceID: "0000-0000-0000",
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/v1/games/1/apps/2/anonymous/sync", bytes.NewBuffer(reqAnomBody))
	router.ServeHTTP(w, req)
	s.Require().Equal(200, w.Code)

	bs := w.Body.Bytes()
	anomResp := &signupAnomResponseTest{}
	err = json.Unmarshal(bs, &anomResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(anomResp.UserID)

	reqUserBody, err := json.Marshal(&signupRealRequestTest{
		GuestID:   anomResp.UserID,
		DeviceID:  "0000-0000-0000",
		Network:   "GOOGLE",
		NetworkID: "google-id",
	})
	s.Require().NoError(err)

	var usersCount int

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/auth/v1/games/1/apps/2/users/sync", bytes.NewBuffer(reqUserBody))
	router.ServeHTTP(w, req)
	s.Require().Equal(200, w.Code)
	bs = w.Body.Bytes()
	realResp := &signupRealResponseTest{}
	err = json.Unmarshal(bs, &realResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(realResp.UserID)
	s.Require().Equal(anomResp.UserID, realResp.GuestID)
	s.Require().Equal(anomResp.UserID, realResp.UserID)
	s.Require().NotEqual(anomResp.JWT, realResp.JWT)
	row := s.PostgresPool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM users")
	row.Scan(&usersCount)
	s.Require().Equal(usersCount, 1)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/auth/v1/games/1/apps/2/users/sync", bytes.NewBuffer(reqUserBody))
	router.ServeHTTP(w, req)
	s.Require().Equal(200, w.Code)
	bs = w.Body.Bytes()
	realResp = &signupRealResponseTest{}
	err = json.Unmarshal(bs, &realResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(realResp.UserID)
	s.Require().Equal(anomResp.UserID, realResp.GuestID)
	s.Require().Equal(anomResp.UserID, realResp.UserID)
	s.Require().NotEqual(anomResp.JWT, realResp.JWT)
	row = s.PostgresPool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM users")
	row.Scan(&usersCount)
	s.Require().Equal(usersCount, 1)

	reqUserBody, err = json.Marshal(&signupRealRequestTest{
		GuestID:   anomResp.UserID,
		DeviceID:  "0000-0000-0000",
		Network:   "FACEBOOK",
		NetworkID: "facebook-id",
	})
	s.Require().NoError(err)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/auth/v1/games/1/apps/2/users/sync", bytes.NewBuffer(reqUserBody))
	router.ServeHTTP(w, req)
	s.Require().Equal(200, w.Code)
	bs = w.Body.Bytes()
	realResp = &signupRealResponseTest{}
	err = json.Unmarshal(bs, &realResp)
	s.Require().NoError(err)
	s.Require().NotEmpty(realResp.UserID)
	s.Require().Equal(anomResp.UserID, realResp.GuestID)
	s.Require().Equal(anomResp.UserID, realResp.UserID)
	s.Require().NotEqual(anomResp.JWT, realResp.JWT)
	row = s.PostgresPool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM users")
	row.Scan(&usersCount)
	s.Require().Equal(usersCount, 2)
}
