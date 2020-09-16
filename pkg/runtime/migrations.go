package runtime

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	// import driver postgres to apply migrations
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"gitlab.com/balconygames/analytics/pkg/postgres"
)

const DefaultDB = "postgres"

func (r *Runtime) WithMigrations(path string, pgConfig postgres.Config) error {
	// only allowed in dev
	if r.action == "db.reset" {
		if !r.spec.Dev() {
			return errors.Errorf("db reset can't be runnable in dev environment", r.spec.Env)
		}

		err := r.drop(pgConfig)
		if err != nil {
			return err
		}
	}

	if r.action == "db.migrate" || r.action == "db.reset" {
		err := r.migrate(path, pgConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) createDB(cfg postgres.Config) error {
	dbName := cfg.Name
	cfg.Name = DefaultDB
	pool, err := pgxpool.Connect(context.Background(), cfg.URL())
	if err != nil {
		return errors.WithStack(err)
	}

	q := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1);`
	var exists bool
	err = pool.QueryRow(context.Background(), q, dbName).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	fmt.Printf("[MIGRATE] CREATING db %s \n", dbName)
	_, err = pool.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE %s;`, pgx.Identifier{dbName}.Sanitize()))
	return err
}

func (r *Runtime) migrate(path string, pgConfig postgres.Config) error {
	fmt.Printf("[MIGRATE] %s to %s \n", path, pgConfig.Name)
	err := r.createDB(pgConfig)
	if err != nil {
		return err
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", path),
		pgConfig.URL(),
	)

	if err != nil {
		return err
	}

	if v := upVersion(); v != 0 {
		fmt.Printf("[MIGRATE] force migrate to %d\n", v)
		m.Force(v)
	}

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r *Runtime) drop(pgConfig postgres.Config) error {
	dbName := pgConfig.Name
	fmt.Printf("[DROP] %s \n", dbName)
	pgConfig.Name = DefaultDB
	pool, err := pgxpool.Connect(context.Background(), pgConfig.URL())
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = pool.Exec(
		context.Background(),
		fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, pgx.Identifier{dbName}.Sanitize()),
	)
	return err
}

func upVersion() int {
	if len(os.Args) < 3 {
		return 0
	}

	i, _ := strconv.Atoi(os.Args[2])
	return i
}
