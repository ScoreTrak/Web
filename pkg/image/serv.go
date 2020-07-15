package image

type Serv interface {
	Delete(id uint64) error
	GetAll() ([]*Image, error)
	GetByID(id uint64) (*Image, error)
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

func (svc *userServ) Delete(id uint64) error { return svc.repo.Delete(id) }

func (svc *userServ) GetAll() ([]*Image, error) { return svc.repo.GetAll() }

func (svc *userServ) GetByID(id uint64) (*Image, error) { return svc.repo.GetByID(id) }

func (svc *userServ) Store(u *Image) error { return svc.repo.Store(u) }

func (svc *userServ) Update(u *Image) error { return svc.repo.Update(u) }
