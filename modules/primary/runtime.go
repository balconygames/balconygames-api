package primary

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.com/balconygames/analytics/modules/primary/internal/db"
	"gitlab.com/balconygames/analytics/modules/primary/internal/handlers"
	"gitlab.com/balconygames/analytics/modules/primary/internal/service"
	"gitlab.com/balconygames/analytics/pkg/logging"
	"gitlab.com/balconygames/analytics/pkg/postgres"
	"gitlab.com/balconygames/analytics/pkg/runtime"
)

type spec struct {
	Env      string          `envconfig:"ENV" required:"True"`
	Postgres postgres.Config `envconfig:"POSTGRES" required:"True"`
}

func New(r *runtime.Runtime) error {
	var s spec
	if err := envconfig.Process("MODULE_PRIMARY", &s); err != nil {
		return err
	}
	return withSpec(r, s)
}

func withSpec(r *runtime.Runtime, s spec) error {
	err := r.WithMigrations("modules/primary/migrations", s.Postgres)
	if err != nil {
		return err
	}

	l, err := logging.ConfigForEnv(s.Env).Build(
		zap.Fields(zap.String("project", "modules/primary")),
	)
	if err != nil {
		return err
	}

	defer l.Sync()
	logger := l.Sugar()

	poolConfig, err := pgxpool.ParseConfig(s.Postgres.URL())
	if err != nil {
		return errors.New("unable to parse database url")
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return errors.Wrap(err, "failed to establish db connection")
	}
	r.WithClosable(pool)

	repo := db.NewPostgresRepository(pool)
	svc := service.NewService(repo, logger)
	h := handlers.New(svc, logger)

	r.WithRoutes(func(r1 chi.Router) {
		// on the client we should give ability to get application info like
		// version name.
		r1.Get("/primary/v1/games/{game_id}/apps/{app_id}", h.GetAppInfo)

		r.WithClientAuth(r1, func(r2 chi.Router) {
		})

		r.WithServerAuth(r1, func(r2 chi.Router) {
			r2.Get("/primary/v1/games", h.ListGames)
			r2.Post("/primary/v1/games", h.CreateGame)
			r2.Put("/primary/v1/games/{game_id}", h.UpdateGame)
			r2.Delete("/primary/v1/games/{game_id}", h.DeleteGame)

			r2.Get("/primary/v1/games/{game_id}/apps", h.ListApps)
			r2.Post("/primary/v1/games/{game_id}/apps", h.CreateApp)
			r2.Put("/primary/v1/games/{game_id}/apps/{app_id}", h.UpdateApp)
			r2.Delete("/primary/v1/games/{game_id}/apps/{app_id}", h.DeleteApp)
		})
	})

	return nil
}
