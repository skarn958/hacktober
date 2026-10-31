package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/secfault-org/hacktober/internal/model/user"
)

type UserRepository interface {
	UserByFingerprint(ctx context.Context, fingerprint string) (*user.User, error)
	UpsertUser(ctx context.Context, fingerprint, name string) (*user.User, error)
	RecordSolve(ctx context.Context, fingerprint, challengeID string) (time.Time, error)
	SolvedBy(ctx context.Context, fingerprint string) (map[string]time.Time, error)
}

type userRepo struct {
	db *sqlx.DB
}

var _ UserRepository = (*userRepo)(nil)

func (r *userRepo) UserByFingerprint(ctx context.Context, fingerprint string) (*user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, `SELECT fingerprint, name FROM users WHERE fingerprint = ?`, fingerprint)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying user %q: %w", fingerprint, err)
	}
	return &u, nil
}

func (r *userRepo) UpsertUser(ctx context.Context, fingerprint, name string) (*user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, `
    	INSERT INTO users (fingerprint, name) VALUES (?, ?)
    	ON CONFLICT(fingerprint) DO UPDATE SET name = excluded.name
    	RETURNING fingerprint, name`, fingerprint, name)
	if err != nil {
		return nil, fmt.Errorf("upserting user %q: %w", fingerprint, err)
	}
	return &u, nil
}

func (r *userRepo) RecordSolve(ctx context.Context, fingerprint, challengeID string) (time.Time, error) {
	var solvedAt time.Time
	err := r.db.GetContext(ctx, &solvedAt, `
              INSERT INTO solves (fingerprint, challenge_id) VALUES (?, ?)
              ON CONFLICT(fingerprint, challenge_id) DO UPDATE SET solved_at = solves.solved_at
              RETURNING solved_at`, fingerprint, challengeID)
	if err != nil {
		return time.Time{}, fmt.Errorf("recording solve %q for %q: %w", challengeID, fingerprint, err)
	}
	return solvedAt, nil
}

func (r *userRepo) SolvedBy(ctx context.Context, fingerprint string) (map[string]time.Time, error) {
	var rows []struct {
		ChallengeID string    `db:"challenge_id"`
		SolvedAt    time.Time `db:"solved_at"`
	}
	err := r.db.SelectContext(ctx, &rows, `SELECT challenge_id, solved_at FROM solves WHERE fingerprint = ?`, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("querying solves for %q: %w", fingerprint, err)
	}

	solved := make(map[string]time.Time, len(rows))
	for _, row := range rows {
		solved[row.ChallengeID] = row.SolvedAt
	}
	return solved, nil
}
