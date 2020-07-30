package image

import "github.com/gofrs/uuid"

type Repo interface {
	Delete(id uuid.UUID) error
	GetAll() ([]*Image, error)
	GetByID(id uuid.UUID) (*Image, error)
	Store(u *Image) error
	Update(u *Image) error
}
