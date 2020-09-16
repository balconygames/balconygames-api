package test_helpers

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/balconygames/analytics/pkg/postgres"
)

type PostgresSuite struct {
	suite.Suite

	Config             postgres.Config
	PostgresPool       *pgxpool.Pool
	PostgresMigrations string
}

func NewDefaultPostgresSuite(t *testing.T, dbName string) PostgresSuite {
	return PostgresSuite{
		PostgresMigrations: MigrationsFolder(t),
		Config: postgres.Config{
			Host:       "localhost:5440",
			Name:       dbName,
			User:       "postgres",
			DisableSSL: true,
		},
	}
}

func (s *PostgresSuite) newPool(ctx context.Context) *pgxpool.Pool {
	rootConfig := postgres.Config{
		User: s.Config.User,
		Pass: s.Config.Pass,
		Host: s.Config.Host,
		Name: "postgres",
	}
	pool, err := pgxpool.Connect(ctx, rootConfig.URL())
	require.NoError(s.T(), err)

	return pool
}

func (s *PostgresSuite) drop(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, "DROP DATABASE IF EXISTS "+s.Config.Name)
	require.NoError(s.T(), err)
}

func (s *PostgresSuite) create(ctx context.Context, pool *pgxpool.Pool) {
	_, err := pool.Exec(ctx, "CREATE DATABASE "+s.Config.Name)
	require.NoError(s.T(), err)
}

func (s *PostgresSuite) SetupTest() {
	ctx := context.Background()
	pool := s.newPool(ctx)
	defer pool.Close()
	s.drop(ctx, pool)
	s.create(ctx, pool)

	m, err := migrate.New(
		s.PostgresMigrations,
		s.Config.URL())
	require.NoError(s.T(), err)

	err = m.Up()
	if _, ok := err.(*os.PathError); !ok {
		require.NoError(s.T(), err)
	}
	m.Close()

	s.PostgresPool, err = pgxpool.Connect(ctx, s.Config.URL())
	require.NoError(s.T(), err)
}

func (s *PostgresSuite) TearDownTest() {
	waiting := sync.WaitGroup{}

	waiting.Add(1)
	go func() {
		done := make(chan bool)
		go func(d chan bool) {
			s.PostgresPool.Close()
			d <- true
		}(done)

		select {
		case <-done:
			// jackc/puddle is hanging on Close()
			// idleConnections != allConnections
			// no reasons for now.
			waiting.Done()
		case <-time.After(300 * time.Millisecond):
			waiting.Done()
			return
		}
	}()

	waiting.Wait()

	ctx := context.Background()
	pool := s.newPool(ctx)
	defer pool.Close()
	s.drop(ctx, pool)
}

func MigrationsFolder(t *testing.T) string {
	pathElements := []string{"migrations"}
	for i := 0; i < 10; i++ {
		path := filepath.Join(pathElements...)
		matches, err := filepath.Glob(path)
		if err != nil {
			t.Fatal(err)
		}

		if len(matches) != 0 {
			return "file://" + path
		}

		pathElements = append([]string{".."}, pathElements...)
	}

	t.Fatalf("migrations folder not found")
	return ""
}
