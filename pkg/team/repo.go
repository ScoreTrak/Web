package team

type Repo interface {
	Delete(id uint64) error
	GetAll() ([]*Team, error)
	GetByID(id uint64) (*Team, error)
	Store(u *Team) error
	Update(u *Team) error
}
