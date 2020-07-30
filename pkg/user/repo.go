package user

import "github.com/gofrs/uuid"

type Repo interface {
	Delete(id uuid.UUID) error
	GetAll() ([]*User, error)
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	Store(u []*User) error
	Update(u *User) error
}
