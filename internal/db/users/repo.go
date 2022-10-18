package users

import (
	"context"
	"database/sql"
	"github.com/prybintsev/validation_cloud/internal/db"

	"github.com/google/uuid"
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
	rowsAffeccted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffeccted == 0 {
		return db.ErrorUserAlreadyExists
	}
	return nil
}
