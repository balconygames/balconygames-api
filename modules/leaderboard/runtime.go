package leaderboard

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/db"
	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/handlers"
	"gitlab.com/balconygames/analytics/modules/leaderboard/internal/service"
	"gitlab.com/balconygames/analytics/pkg/geo"
	"gitlab.com/balconygames/analytics/pkg/logging"
	"gitlab.com/balconygames/analytics/pkg/postgres"
	redisconf "gitlab.com/balconygames/analytics/pkg/redis"
	"gitlab.com/balconygames/analytics/pkg/runtime"
)

type spec struct {
	Env       string           `envconfig:"ENV" required:"True"`
	Postgres  postgres.Config  `envconfig:"POSTGRES" required:"True"`
	Redis     redisconf.Config `envconfig:"REDIS" required:"True"`
	AES256Key string           `envconfig:"AES256_KEY" required:"True"`
}

func New(r *runtime.Runtime) error {
	var s spec
	if err := envconfig.Process("MODULE_LEADERBOARD", &s); err != nil {
		return err
	}
	return withSpec(r, s)
}

func withSpec(r *runtime.Runtime, s spec) error {
	err := r.WithMigrations("modules/leaderboard/migrations", s.Postgres)
	if err != nil {
		panic(err)
	}

	l, err := logging.ConfigForEnv(s.Env).Build(
		zap.Fields(zap.String("project", "modules/leaderboard")),
	)
	if err != nil {
		return err
	}

	defer l.Sync()
	logger := l.Sugar()

	logger.Debug("postgres connection string", s.Postgres.URL())

	poolConfig, err := pgxpool.ParseConfig(s.Postgres.URL())
	if err != nil {
		return errors.New("unable to parse database url")
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return errors.Wrap(err, "failed to establish db connection")
	}
	r.WithClosable(pool)

	geoResolver := geo.New()
	r.WithClosable(geoResolver)

	repo := db.NewPostgresRepository(pool)

	logger.Debug("connect to redis ", "host", s.Redis.Host, " port ", s.Redis.Port)

	redisConn := redis.NewClient(s.Redis.Options())
	redisRepo := db.NewRedisRepository(redisConn, logger)

	svc := service.NewService(repo, redisRepo, geoResolver, logger)
	h := handlers.New(svc, logger)

	r.WithRoutes(func(r1 chi.Router) {
		r.WithClientAuth(r1, func(r2 chi.Router) {
			r2.Post("/leaderboard/v1/scores", h.CreateScores)
			r2.Post("/leaderboard/v1/scores/list", h.ListScores)
		})

		// Server API would be used by dashboard layer
		r.WithServerAuth(r1, func(r2 chi.Router) {
			r1.Post("/leaderboard/v1/leaderboards", h.CreateLeaderboards)
			r1.Get("/leaderboard/v1/leaderboards/list", h.ListLeaderboards)
		})
	})

	return nil
}
