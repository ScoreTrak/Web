package repo

import (
	"github.com/ScoreTrak/Web/pkg/di"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
)

func NewStore() Store {
	var teamRepo team.Repo
	di.Invoke(func(re team.Repo) {
		teamRepo = re
	})
	var userRepo user.Repo
	di.Invoke(func(re user.Repo) {
		userRepo = re
	})
	var policyRepo policy.Repo
	di.Invoke(func(re policy.Repo) {
		policyRepo = re
	})
	return Store{
		Policy: policyRepo,
		User:   userRepo,
		Team:   teamRepo,
	}
}

type Store struct {
	Team   team.Repo
	User   user.Repo
	Policy policy.Repo
}
