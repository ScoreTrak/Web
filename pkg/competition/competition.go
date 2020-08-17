package competition

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/competition"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
)

type Web struct {
	Competition *competition.Competition
	Teams       []*team.Team
	Users       []*user.User
	Policy      *policy.Policy
}
