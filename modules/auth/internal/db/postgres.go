package db

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.com/balconygames/analytics/modules/auth/internal/models"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const NamePattern = "Guest-"

// AnomSync would create or return user_id from database
func (r PostgresRepository) AnomSync(ctx context.Context, user *models.User) error {
	var err error

	query := `
		INSERT INTO
			anonymouses (
				user_id
				, game_id
				, app_id
				, device_id
				, name
			)
			VALUES (
				$1
				, $2
				, $3
				, $4
				, $5
			)
		ON CONFLICT (
			user_id
			, game_id
			, app_id
		)
		DO UPDATE SET updated_at = NOW()
		RETURNING user_id
	`

	// should set uuid for user_id column
	// only for new records.
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil
	}

	// TODO: allow to pass pattern or custom name on guest sign in
	if user.Name == "" {
		if len(user.DeviceID) > 5 {
			user.Name = NamePattern + user.DeviceID[0:5]
		} else {
			// in case if we have weird device id, we should
			// use 5 letters of guid
			user.Name = NamePattern + uuid[0:5]
		}
	}

	row := r.pool.QueryRow(ctx, query, uuid, user.Scope.GameID, user.AppID, user.DeviceID, user.Name)
	err = row.Scan(&user.Scope.UserID)
	if err != nil {
		return err
	}

	return nil
}

// AnomSync would create or return user_id from database
func (r PostgresRepository) UsersSync(ctx context.Context, user *models.User) error {
	var err error

	query := `
		select
			user_id
			, name
			, guest_id
		from users
		where
			network_name=$1
			and network_id=$2
			and game_id=$3
			and app_id=$4
		limit 1
	`
	row := r.pool.QueryRow(ctx, query, user.Network,
		user.NetworkID, user.Scope.GameID, user.AppID)

	err = row.Scan(&user.Scope.UserID, &user.Name, &user.GuestID)
	if err == nil {
		// if we found user, we should respond with
		// settings.
		return nil
	}

	query = `
		INSERT INTO
			users (
				user_id
				, guest_id
				, device_id

				, game_id
				, app_id

				, network_name
				, network_id

				, email
				, name
			)
			VALUES (
				$1
				, $2
				, $3
				, $4
				, $5
				, $6
				, $7
				, $8
				, $9
			)
		ON CONFLICT (
			user_id
			, network_name
			, network_id
			, game_id
			, app_id
		)
		DO UPDATE SET updated_at = NOW()
		RETURNING user_id
	`

	if user.UserID == "" {
		// use guest id as primary id because of using
		// using only once anom sign in to generate
		// player id.
		user.UserID = user.GuestID
	}

	if user.Name == "" {
		if len(user.DeviceID) > 5 {
			user.Name = NamePattern + user.DeviceID[0:5]
		} else {
			// in case if we have weird device id, we should
			// use 5 letters of guid
			user.Name = NamePattern + user.UserID[0:5]
		}
	}

	row = r.pool.QueryRow(ctx, query, user.UserID, user.GuestID,
		user.DeviceID, user.Scope.GameID, user.AppID,
		user.Network, user.NetworkID, user.Email, user.Name)
	err = row.Scan(&user.Scope.UserID)
	if err != nil {
		return err
	}

	return nil
}
