package track

import "github.com/secfault-org/hacktober/internal/model/challenge"

type Track struct {
	Id          string
	Title       string
	Description string
	Challenges  []challenge.Challenge
}
