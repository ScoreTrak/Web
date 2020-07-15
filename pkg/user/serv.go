package user

type Serv interface {
	Delete(id uint64) error
	GetAll() ([]*User, error)
	GetByID(id uint64) (*User, error)
	GetByUsername(username string) (*User, error)
	Store(u *User) error
	Update(u *User) error
}

type userServ struct {
	repo Repo
}

func NewUserServ(repo Repo) Serv {
	return &userServ{
		repo: repo,
	}
}

func (svc *userServ) Delete(id uint64) error { return svc.repo.Delete(id) }

func (svc *userServ) GetAll() ([]*User, error) { return svc.repo.GetAll() }

func (svc *userServ) GetByID(id uint64) (*User, error) { return svc.repo.GetByID(id) }

func (svc *userServ) GetByUsername(username string) (*User, error) {
	return svc.repo.GetByUsername(username)
}

func (svc *userServ) Store(u *User) error { return svc.repo.Store(u) }

func (svc *userServ) Update(u *User) error { return svc.repo.Update(u) }
