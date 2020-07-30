package team

import "github.com/gofrs/uuid"

type Repo interface {
	Delete(id uuid.UUID) error
	GetAll() ([]*Team, error)
	GetByID(id uuid.UUID) (*Team, error)
	Store(u []*Team) error
	Update(u *Team) error
}
