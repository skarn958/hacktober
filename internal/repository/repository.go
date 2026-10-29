package repository

import (
	"context"

	"github.com/secfault-org/hacktober/internal/config"
)

type Repository interface {
	ChallengeRepository
}

type datastore struct {
	ctx context.Context

	*challengeRepo
}

func NewRepository(ctx context.Context, cfg *config.Config) Repository {
	return &datastore{
		ctx: ctx,
		challengeRepo: &challengeRepo{
			basedir: cfg.Challenge.BaseDir,
			tracks:  cfg.Tracks,
		},
	}
}
