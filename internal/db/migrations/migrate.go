package migrations

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	_ "github.com/golang-migrate/migrate/source/file"
)

func MigrateUp(db *sql.DB, dbName string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations/scripts",
		dbName, driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("Migrations: no change")
			return nil
		}
		return err
	}

	log.Info("Migrations have completed successfully")

	return nil
}
