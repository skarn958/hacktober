package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/secfault-org/hacktober/internal/config"
)

type Repository interface {
	ChallengeRepository
	UserRepository
}

type datastore struct {
	*challengeRepo
	*userRepo
}

func NewRepository(_ context.Context, cfg *config.Config, db *sqlx.DB) Repository {
	return &datastore{
		challengeRepo: &challengeRepo{
			basedir: cfg.Challenge.BaseDir,
			tracks:  cfg.Tracks,
		},
		userRepo: &userRepo{db: db},
	}
}
