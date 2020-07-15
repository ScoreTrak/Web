package user

type Repo interface {
	Delete(id uint64) error
	GetAll() ([]*User, error)
	GetByID(id uint64) (*User, error)
	GetByUsername(username string) (*User, error)
	Store(u *User) error
	Update(u *User) error
}
