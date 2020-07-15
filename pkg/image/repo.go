package image

type Repo interface {
	Delete(id uint64) error
	GetAll() ([]*Image, error)
	GetByID(id uint64) (*Image, error)
	Store(u *Image) error
	Update(u *Image) error
}
