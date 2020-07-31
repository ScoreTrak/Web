package policy

type Repo interface {
	Get() (*Policy, error)
	Update(u *Policy) error
}
