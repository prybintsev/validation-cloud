package users

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/prybintsev/validation_cloud/internal/db"
)

type Users struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *Users {
	return &Users{db: db}
}

func (u *Users) CreateUser(ctx context.Context, userName, passwordHash string) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	res, err := u.db.ExecContext(ctx, "INSERT OR IGNORE INTO user (ID, Username, PasswordHash) VALUES  (?, ?, ?)",
		id.String(), userName, passwordHash)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return db.ErrorUserAlreadyExists
	}
	return nil
}

func (u *Users) GetPasswordHashByUsername(ctx context.Context, userName string) (string, error) {
	rows, err := u.db.QueryContext(ctx, "SELECT PasswordHash FROM user WHERE Username = ?", userName)
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.WithError(err).Error("failed to close query response")
		}
	}()

	var passHash string
	if !rows.Next() {
		return "", db.ErrorUserNotFound
	}

	err = rows.Scan(&passHash)
	if err != nil {
		return "", err
	}

	return passHash, nil
}
