package policy

type Serv interface {
	Get() (*Policy, error)
	Update(u *Policy) error
}

type policyServ struct {
	repo Repo
}

func NewPolicyServ(repo Repo) Serv {
	return &policyServ{
		repo: repo,
	}
}

func (svc *policyServ) Get() (*Policy, error) { return svc.repo.Get() }

func (svc *policyServ) Update(u *Policy) error { return svc.repo.Update(u) }
