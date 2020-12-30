package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
	test_helpers "gitlab.com/balconygames/analytics/pkg/test_helpers"
	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type serviceSuite struct {
	test_helpers.PostgresSuite
	logger *zap.SugaredLogger
}

const deviceID = "device-1"

func TestEntrypointSuite(t *testing.T) {
	handler := &serviceSuite{
		PostgresSuite: test_helpers.NewDefaultPostgresSuite(t, "test_auth"),
		logger:        zaptest.NewLogger(t).Sugar(),
	}

	suite.Run(t, handler)
}

func (s *serviceSuite) TestAnomSync() {
	repo := NewPostgresRepository(s.PostgresPool)

	_ = repo.AnomSync(context.Background(), &models.User{
		Scope: sharedmodels.Scope{
			AppID:  appID,
			GameID: gameID,
		},
		DeviceID: "device-1",
	})

	var id, deviceID string
	row := s.PostgresPool.QueryRow(context.Background(), `SELECT user_id, device_id FROM anonymouses LIMIT 1`)
	_ = row.Scan(&id, &deviceID)
	s.Require().NotEmpty(id)
	s.Require().NotEmpty(deviceID)
	s.Require().Equal("device-1", deviceID)

	var anomCount int64
	countRow := s.PostgresPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM anonymouses LIMIT 1`)
	_ = row.Scan(&countRow)
	s.Require().Equal(1, anomCount)
}

func (s *serviceSuite) TestAnomSyncNoDuplicates() {
	repo := NewPostgresRepository(s.PostgresPool)

	user1 := &models.User{
		Scope: sharedmodels.Scope{
			AppID:  appID,
			GameID: gameID,
		},
		DeviceID: deviceID,
	}
	_ = repo.AnomSync(context.Background(), user1)

	user2 := &models.User{
		Scope: sharedmodels.Scope{
			AppID:  appID,
			GameID: gameID,
		},
		DeviceID: deviceID,
	}
	_ = repo.AnomSync(context.Background(), user2)

	var id, deviceID string
	row := s.PostgresPool.QueryRow(context.Background(), `SELECT user_id, device_id FROM anonymouses LIMIT 1`)
	_ = row.Scan(&id, &deviceID)
	s.Require().NotEmpty(id)
	s.Require().NotEmpty(user1.UserID)
	s.Require().Equal(user1.UserID, id)
	s.Require().Equal(deviceID, deviceID)

	var anomCount int64
	countRow := s.PostgresPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM anonymouses LIMIT 1`)
	_ = row.Scan(&countRow)
	s.Require().Equal(1, anomCount)
}

func (s *serviceSuite) TestAnomSyncSameUsers() {
	var row pgx.Row
	var id1, id2 string

	repo := NewPostgresRepository(s.PostgresPool)

	_ = repo.AnomSync(context.Background(), &models.User{
		Scope: sharedmodels.Scope{
			AppID:  appID,
			GameID: gameID,
		},
		DeviceID: deviceID,
	})
	row = s.PostgresPool.QueryRow(context.Background(), `SELECT user_id FROM anonymouses LIMIT 1`)
	_ = row.Scan(&id1)

	_ = repo.AnomSync(context.Background(), &models.User{
		Scope: sharedmodels.Scope{
			AppID:  appID,
			GameID: gameID,
		},
		DeviceID: deviceID,
	})
	row = s.PostgresPool.QueryRow(context.Background(), `SELECT user_id FROM anonymouses LIMIT 1`)
	_ = row.Scan(&id2)

	s.Require().Equal(id1, id2)
}
