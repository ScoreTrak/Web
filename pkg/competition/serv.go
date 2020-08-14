package competition

import (
	"errors"
	"github.com/ScoreTrak/Web/pkg/di/repo"
	"github.com/jackc/pgconn"
)

type Serv interface {
	LoadCompetition(*Web) error
	FetchCompetition() (*Web, error)
}

type webServ struct {
	Store repo.Store
}

func NewCompetitionServ(str repo.Store) Serv {
	return &webServ{
		Store: str,
	}
}

func (svc *webServ) LoadCompetition(c *Web) error {
	var errAgr []error
	errAgr = append(errAgr, svc.Store.Team.Store(c.Teams))
	errAgr = append(errAgr, svc.Store.User.Store(c.Users))
	errAgr = append(errAgr, svc.Store.Policy.Update(c.Policy))
	errStr := ""
	for i, _ := range errAgr {
		if errAgr[i] != nil {
			serr, ok := errAgr[i].(*pgconn.PgError)
			if !ok || serr.Code != "23505" {
				errStr += errAgr[i].Error() + "\n"
			}
		}
	}
	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}

func (svc *webServ) FetchCompetition() (*Web, error) {
	policy, err := svc.Store.Policy.Get()
	if err != nil {
		return nil, err
	}
	teams, err := svc.Store.Team.GetAll()
	if err != nil {
		return nil, err
	}
	users, err := svc.Store.User.GetAll()
	if err != nil {
		return nil, err
	}

	return &Web{Teams: teams, Policy: policy, Users: users}, nil
}
