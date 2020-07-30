package image

import "github.com/gofrs/uuid"

type Serv interface {
	Delete(id uuid.UUID) error
	GetAll() ([]*Image, error)
	GetByID(id uuid.UUID) (*Image, error)
	Store(u *Image) error
	Update(u *Image) error
}

type userServ struct {
	repo Repo
}

func NewImageServ(repo Repo) Serv {
	return &userServ{
		repo: repo,
	}
}

func (svc *userServ) Delete(id uuid.UUID) error { return svc.repo.Delete(id) }

func (svc *userServ) GetAll() ([]*Image, error) { return svc.repo.GetAll() }

func (svc *userServ) GetByID(id uuid.UUID) (*Image, error) { return svc.repo.GetByID(id) }

func (svc *userServ) Store(u *Image) error { return svc.repo.Store(u) }

func (svc *userServ) Update(u *Image) error { return svc.repo.Update(u) }
