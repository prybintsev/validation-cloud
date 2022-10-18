package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/prybintsev/validation_cloud/internal/db"
	"github.com/prybintsev/validation_cloud/internal/db/migrations"
)

const (
	DBFile = "validation-cloud.db"
	DBDir  = "data"
)

func ConnectAndMigrate(ctx context.Context) (*sql.DB, error) {
	err := os.MkdirAll(DBDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	dbFilePath := fmt.Sprintf("%s/%s", DBDir, DBFile)
	_, err = os.Stat(fmt.Sprintf("%s/%s", DBDir, DBFile))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var file *os.File
			file, err = os.Create(dbFilePath)
			if err != nil {
				log.WithError(err).WithField("dbFile", dbFilePath).Error("Failed to crete db file")
				return nil, err
			}
			defer file.Close()
			log.WithField("dbFile", dbFilePath).Info("Created not open sqlite data file")
		} else {
			log.WithError(err).WithField("dbFile", dbFilePath).Error("Could not open sqlite data file")
			return nil, err
		}
	}

	dbCon, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.WithError(err).WithField("dbFile", dbFilePath).Error("Could not open sqlite database")
		return nil, err
	}
	err = migrations.MigrateUp(dbCon, db.DBName)
	if err != nil {
		log.WithError(err).WithField("dbFile", dbFilePath).Error("Failed to perform migrations")
		return nil, err
	}

	return dbCon, nil
}
