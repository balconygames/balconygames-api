package auth

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"gitlab.com/balconygames/analytics/pkg/runtime"
	th "gitlab.com/balconygames/analytics/pkg/test_helpers"
)

type serviceSuite struct {
	suite.Suite

	th.PostgresSuite
	th.RedisSuite

	logger *zap.SugaredLogger
}

func TestEntrypointSuite(t *testing.T) {
	handler := &serviceSuite{
		PostgresSuite: th.NewDefaultPostgresSuite(t, "test_auth"),
		RedisSuite:    th.NewDefaultRedisSuite(t, "test_auth"),
		logger:        zaptest.NewLogger(t).Sugar(),
	}

	suite.Run(t, handler)
}

func (s *serviceSuite) TestAuth() {
	sp := runtime.Spec{
		Env: "test",
	}
	r := runtime.New("web", sp)
	aesKey := "<aes>"
	err := withSpec(r, spec{
		Env:       "test",
		Postgres:  s.PostgresSuite.Config,
		AES256Key: aesKey,
	})
	s.Require().NoError(err)

	go func() {
		err = r.Run()
		s.Require().NoError(err)
	}()
	r.Wait()
	defer r.Close()
}
