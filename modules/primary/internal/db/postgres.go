package db

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	sharedmodels "gitlab.com/balconygames/analytics/shared/models"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r PostgresRepository) CreateGame(ctx context.Context, game *sharedmodels.Game) error {
	var err error

	query := `
		INSERT INTO
			games (
				game_id
				, name
			)
			VALUES (
				$1
				, $2
			)
		ON CONFLICT (
			game_id
		)
		DO UPDATE SET
			updated_at = NOW()
			, NAME = $2
		RETURNING game_id
	`

	if game.GameID == "" {
		// should set uuid for user_id column
		// only for new records.
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return nil
		}

		game.GameID = uuid
	}

	row := r.pool.QueryRow(ctx, query, game.GameID, game.Name)
	err = row.Scan(&game.GameID)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresRepository) CreateApp(ctx context.Context, app *sharedmodels.App) error {
	var err error

	query := `
		INSERT INTO
			apps (
				game_id
				, app_id
				, platform
				, market
				, version
				, force_update_enabled
				, device_type
			)
			VALUES (
				$1
				, $2
				, $3
				, $4
				, $5
				, $6
				, $7
			)
		ON CONFLICT (
			game_id
			, app_id
		)
		DO UPDATE SET
			updated_at = NOW()
			, platform=$3
			, market = $4
			, device_type = $7
		RETURNING app_id
	`

	if app.AppID == "" {
		// should set uuid for user_id column
		// only for new records.
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return nil
		}

		app.AppID = uuid
	}

	row := r.pool.QueryRow(ctx, query, app.GameID,
		app.AppID, app.Platform, app.Market,
		app.Version, app.ForceUpdateEnabled,
		app.DeviceType)
	err = row.Scan(&app.AppID)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresRepository) GetAppInfo(ctx context.Context, app *sharedmodels.App) error {
	var err error

	query := `
		SELECT
			game_id
			, app_id
			, platform
			, market
			, device_type
			, version
			, created_at
			, updated_at
			, force_update_enabled
		FROM apps
		WHERE
			game_id=$1
			AND app_id=$2
	`

	row := r.pool.QueryRow(ctx, query, app.GameID, app.AppID)
	err = row.Scan(&app.GameID, &app.AppID, &app.Platform, &app.Market, &app.DeviceType,
		&app.Version, &app.CreatedAt, &app.UpdatedAt, &app.ForceUpdateEnabled)
	if err != nil {
		return err
	}

	return nil
}
