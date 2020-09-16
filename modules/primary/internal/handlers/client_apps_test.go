package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"gitlab.com/balconygames/analytics/modules/primary/internal/db"
	"gitlab.com/balconygames/analytics/modules/primary/internal/service"
	"gitlab.com/balconygames/analytics/pkg/runtime"
	test_helpers "gitlab.com/balconygames/analytics/pkg/test_helpers"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type serviceSuite struct {
	test_helpers.PostgresSuite

	logger *zap.SugaredLogger
}

func TestEntrypointSuite(t *testing.T) {
	handler := &serviceSuite{
		PostgresSuite: test_helpers.NewDefaultPostgresSuite(t, "test_auth"),
		logger:        zaptest.NewLogger(t).Sugar(),
	}

	suite.Run(t, handler)
}

type signupResponseTest struct {
	sharedmodels.App
}

func (s *serviceSuite) TestRequestHandlerToVerifyToken() {
	var err error

	repo := db.NewPostgresRepository(s.PostgresPool)
	svc := service.NewService(repo, s.logger)

	h := New(svc, s.logger)

	spec := runtime.Spec{
		Env:             "test",
		JWTClientSecret: "<secret>",
	}
	r := runtime.New("api", spec)
	s.Require().NoError(err)

	var router chi.Router

	game := &sharedmodels.Game{Name: "TEST"}
	err = repo.CreateGame(context.Background(), game)
	s.Require().NoError(err)

	app := &sharedmodels.App{
		Platform:   sharedmodels.AndroidPlatform,
		Market:     sharedmodels.HuaweiMarket,
		GameID:     game.GameID,
		DeviceType: sharedmodels.PadDeviceType,
		Version:    "1.0.0"}
	err = repo.CreateApp(context.Background(), app)
	s.Require().NoError(err)

	r.WithRoutes(func(r1 chi.Router) {
		router = r1

		r1.Get("/primary/v1/games/{game_id}/apps/{app_id}", h.GetAppInfo)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("/primary/v1/games/%s/apps/%s", game.GameID, app.AppID),
		nil)
	router.ServeHTTP(w, req)
	s.Require().Equal(200, w.Code)

	respBody := w.Body.Bytes()
	resp := &signupResponseTest{}
	err = json.Unmarshal(respBody, &resp)
	s.Require().NoError(err)
	s.Require().NotEmpty(resp.AppID)
}
