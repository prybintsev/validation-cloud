package samples

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/prybintsev/validation_cloud/internal/db"
	"time"

	log "github.com/sirupsen/logrus"
)

type Samples struct {
	db *sql.DB
}

func NewSamplesRepo(db *sql.DB) Samples {
	return Samples{db: db}
}

type sample struct {
	height    int64
	createdAt time.Time
}

func (s Samples) GetAverageGrowth(ctx context.Context) (float64, error) {
	now := time.Now().UTC()
	oneHourAgo := now.Add(-1 * time.Hour)
	lastSample, err := s.getLastSample(ctx, oneHourAgo, now)
	if err != nil {
		return 0, err
	}
	firstSample, err := s.getFirstSample(ctx, oneHourAgo, now)
	if err != nil {
		return 0, err
	}

	diffInMinutes := lastSample.createdAt.Sub(firstSample.createdAt).Minutes()
	heightDiff := lastSample.height - firstSample.height
	if diffInMinutes == 0 {
		return 0, errors.New("no time difference between two samples")
	}
	return float64(heightDiff) / diffInMinutes, nil
}

func (s Samples) getLastSample(ctx context.Context, from, to time.Time) (sample, error) {
	query := "SELECT BlockchainHeight, CreatedAt FROM sample WHERE CreatedAt BETWEEN ? AND ? ORDER BY CreatedAt DESC LIMIT 1"
	return s.getSample(ctx, query, from, to)
}

func (s Samples) getFirstSample(ctx context.Context, from, to time.Time) (sample, error) {
	query := "SELECT BlockchainHeight, CreatedAt FROM sample WHERE CreatedAt BETWEEN ? AND ? ORDER BY CreatedAt ASC LIMIT 1"
	return s.getSample(ctx, query, from, to)
}

func (s Samples) getSample(ctx context.Context, query string, from, to time.Time) (sample, error) {
	rows, err := s.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return sample{}, err
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			log.WithError(err).Error("failed to close query response")
		}
	}()

	var height int64
	var createdAt time.Time
	if !rows.Next() {
		return sample{}, errors.New("no samples found")
	}

	err = rows.Scan(&height, &createdAt)
	if err != nil {
		return sample{}, err
	}

	return sample{height: height, createdAt: createdAt}, nil
}

func (s Samples) InsertSample(ctx context.Context, height uint64, createdAt time.Time) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	res, err := s.db.ExecContext(ctx, "INSERT INTO sample (ID, BlockchainHeight, CreatedAt) VALUES  (?, ?, ?)",
		id.String(), height, createdAt)
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
