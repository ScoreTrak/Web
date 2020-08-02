package image

import "github.com/gofrs/uuid"

type Serv interface {
	Delete(id uuid.UUID) error
	GetAll() ([]*Image, error)
	GetByID(id uuid.UUID) (*Image, error)
	Store(u *Image) error
	Update(u *Image) error
}

type imageServ struct {
	repo Repo
}

func NewImageServ(repo Repo) Serv {
	return &imageServ{
		repo: repo,
	}
}

func (svc *imageServ) Delete(id uuid.UUID) error { return svc.repo.Delete(id) }

func (svc *imageServ) GetAll() ([]*Image, error) { return svc.repo.GetAll() }

func (svc *imageServ) GetByID(id uuid.UUID) (*Image, error) { return svc.repo.GetByID(id) }

func (svc *imageServ) Store(u *Image) error { return svc.repo.Store(u) }

func (svc *imageServ) Update(u *Image) error { return svc.repo.Update(u) }
