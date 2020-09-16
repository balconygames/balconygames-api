package auth

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/auth/internal/db"
	"gitlab.com/balconygames/analytics/modules/auth/internal/handlers"
	"gitlab.com/balconygames/analytics/modules/auth/internal/service"
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
	if err := envconfig.Process("MODULE_AUTH", &s); err != nil {
		return err
	}
	return withSpec(r, s)
}

func withSpec(r *runtime.Runtime, s spec) error {
	err := r.WithMigrations("modules/auth/migrations", s.Postgres)
	if err != nil {
		panic(err)
	}

	l, err := logging.ConfigForEnv(s.Env).Build(
		zap.Fields(zap.String("project", "modules/auth")),
	)
	if err != nil {
		return err
	}

	defer l.Sync()
	logger := l.Sugar()

	logger.Debug("postgres connection string", s.Postgres.URL())

	pool, err := pgxpool.Connect(context.Background(), s.Postgres.URL())
	if err != nil {
		return errors.Wrap(err, "failed to establish db connection")
	}

	redisConn := redis.NewClient(s.Redis.Options())
	redisRepo := db.NewRedisRepository(redisConn, logger)
	repo := db.NewPostgresRepository(pool)
	svc := service.NewService(repo, redisRepo, logger)
	h := handlers.New(svc, logger)

	r.WithClosable(pool)
	r.WithRoutes(func(r1 chi.Router) {
		// ==== BEGIN CLIENT routes
		// pass token signer instance in context
		r.WithClientTokenSigner(r1, func(r2 chi.Router) {
			r2.Post("/auth/v1/games/{game_id}/apps/{app_id}/anonymous/sync", h.SyncHandler)
			r2.Post("/auth/v1/games/{game_id}/apps/{app_id}/users/sync", h.UsersSyncHandler)

			r2.Group(func(i chi.Router) {
				i.Use(h.FacebookMiddleware)
				i.Post("/auth/v1/games/{game_id}/apps/{app_id}/facebook/login", h.FacebookLogin)
				i.Post("/auth/v1/games/{game_id}/apps/{app_id}/facebook/callback", h.FacebookCallback)
			})

			r2.Group(func(i chi.Router) {
				i.Use(h.GoogleMiddleware)
				i.Post("/auth/v1/games/{game_id}/apps/{app_id}/google/login", h.GoogleLogin)
				i.Post("/auth/v1/games/{game_id}/apps/{app_id}/google/callback", h.GoogleCallback)
			})

			r2.Group(func(i chi.Router) {
				i.Use(h.TwitterMiddleware)
				i.Post("/auth/v1/games/{game_id}/apps/{app_id}/twitter/login", h.TwitterLogin)
				i.Post("/auth/v1/games/{game_id}/apps/{app_id}/twitter/callback", h.TwitterCallback)
			})
		})

		r.WithClientAuth(r1, func(r2 chi.Router) {
			r2.Post("/auth/v1/properties/list", h.GetPropertiesHandler)
			r2.Put("/auth/v1/properties", h.SetPropertiesHandler)
		})
		// ==== END CLIENT routes

		// ==== BEGIN SERVER routes
		r.WithServerTokenSigner(r1, func(r2 chi.Router) {
			// Example: Add facebook sign in
			//
			// r2.Group(func(i chi.Router) {
			// 	i.Use(h.FacebookMiddleware)
			// 	i.Post("/auth/v1/signin/facebook", h.ServerFacebookLogin)
			// 	i.Post("/auth/v1/signin/facebook/callback", h.ServerFacebookCallback)
			// })

			r2.Post("/auth/v1/signin", h.ServerSignin)
			r2.Post("/auth/v1/signup", h.ServerSignup)
		})
		r1.Delete("/auth/v1/signout", h.ServerSignout)

		r.WithServerAuth(r1, func(r2 chi.Router) {
			// TODO: add ability to view properties stored per player
		})
		// ==== END SERVER routes
	})

	return nil
}
