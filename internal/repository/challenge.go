package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"charm.land/log/v2"
	"github.com/secfault-org/hacktober/internal/config"
	"github.com/secfault-org/hacktober/internal/model/challenge"
	"github.com/secfault-org/hacktober/internal/model/track"
)

type ChallengeRepository interface {
	GetAllChallenges(ctx context.Context) ([]challenge.Challenge, error)
	GetAllTracks(ctx context.Context) ([]track.Track, error)
}

type challengeRepo struct {
	basedir string
	tracks  []config.TrackConfig
}

type challengeInfo struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	ShortDesc      string `json:"shortDescription"`
	DescFile       string `json:"descriptionFile"`
	Difficulty     uint8  `json:"difficulty"`
	ReleaseDate    string `json:"releaseDate"`
	ContainerImage string `json:"containerImage"`
	ChallengeFile  string `json:"challengeFile"`
	Track          string `json:"track"`
}

var _ ChallengeRepository = (*challengeRepo)(nil)

func (f *challengeRepo) GetAllChallenges(ctx context.Context) ([]challenge.Challenge, error) {
	var challenges []challenge.Challenge

	err := filepath.WalkDir(f.basedir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "challenge.json" {
			return nil
		}
		chall, err := readChallenge(path)
		if err != nil {
			return fmt.Errorf("reading challenge file: %w", err)
		}
		challenges = append(challenges, chall)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return challenges, nil
}

func readChallenge(challengeJsonPath string) (challenge.Challenge, error) {
	challengeDir := filepath.Dir(challengeJsonPath)
	jsonFile, err := os.ReadFile(challengeJsonPath)
	if err != nil {
		return challenge.Challenge{}, err
	}

	var challengeJson challengeInfo
	err = json.Unmarshal(jsonFile, &challengeJson)
	if err != nil {
		return challenge.Challenge{}, err
	}

	releaseDate, err := time.Parse(time.DateOnly, challengeJson.ReleaseDate)
	if err != nil {
		return challenge.Challenge{}, err
	}

	text, err := os.ReadFile(path.Join(challengeDir, challengeJson.DescFile))
	if err != nil {
		return challenge.Challenge{}, err
	}

	challengeFile, err := filepath.Abs(path.Join(challengeDir, challengeJson.ChallengeFile))
	if err != nil {
		return challenge.Challenge{}, err
	}

	return challenge.Challenge{
		Id:                challengeJson.Id,
		Name:              challengeJson.Title,
		Description:       challengeJson.ShortDesc,
		ChallengeMarkdown: string(text),
		ReleaseDate:       releaseDate,
		ContainerImage:    challengeJson.ContainerImage,
		ChallengeFile:     challengeFile,
		Track:             challengeJson.Track,
	}, nil
}

func (f *challengeRepo) GetAllTracks(ctx context.Context) ([]track.Track, error) {
	challenges, err := f.GetAllChallenges(ctx)
	if err != nil {
		return nil, err
	}

	byTrack := make(map[string][]challenge.Challenge)
	for _, ch := range challenges {
		if ch.Track == "" {
			log.Warn("Challenge has no track, ignoring it", "ChallengeId", ch.Id)
			continue
		}
		byTrack[ch.Track] = append(byTrack[ch.Track], ch)
	}

	var tracks []track.Track
	for _, t := range f.tracks {
		challs := byTrack[t.Id]
		delete(byTrack, t.Id)
		tracks = append(
			tracks,
			track.Track{
				Id:          t.Id,
				Title:       t.Title,
				Description: t.Description,
				Challenges:  challs,
			})
	}
	for id := range byTrack {
		log.Warn("challenge track not listed in config", "TrackID", id)
	}

	return tracks, nil

}
