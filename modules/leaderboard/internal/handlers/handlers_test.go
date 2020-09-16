package handlers

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	test_helpers "gitlab.com/balconygames/analytics/pkg/test_helpers"
)

type serviceSuite struct {
	test_helpers.PostgresSuite
	logger *zap.SugaredLogger
}

func TestEntrypointSuite(t *testing.T) {
	handler := &serviceSuite{
		PostgresSuite: test_helpers.NewDefaultPostgresSuite(t, "test_leaderboard"),
		logger:        zaptest.NewLogger(t).Sugar(),
	}

	suite.Run(t, handler)
}

func (s *serviceSuite) TestRequestHandlerToVerifyToken() {
}
