package policy

import (
	"sync"
	"time"
)

type Client struct {
	policy *Policy
	mu     sync.RWMutex
	repo   Repo
	cnf    ClientConfig
}

func NewPolicyClient(policy *Policy, repo Repo, cnf ClientConfig) *Client {
	return &Client{policy: policy, mu: sync.RWMutex{}, repo: repo, cnf: cnf}
}

func (a *Client) LazyPolicyLoader() {
	t := time.Now()
	rounded := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	time.Sleep(time.Until(rounded))
	for {
		tempPolicy, _ := a.repo.Get()
		a.mu.Lock()
		a.policy.AllowUnauthenticatedUsers = tempPolicy.AllowUnauthenticatedUsers
		a.mu.Unlock()
		rounded = rounded.Add(time.Second * time.Duration(a.cnf.PolicyRefreshSeconds))
		time.Sleep(time.Until(rounded))
	}
}

func (a *Client) GetPolicy() *Policy {
	a.mu.RLock()
	p := &Policy{AllowUnauthenticatedUsers: a.policy.AllowUnauthenticatedUsers, ShowPoints: a.policy.ShowPoints}
	a.mu.RUnlock()
	return p
}

type ClientConfig struct {
	PolicyRefreshSeconds uint `default:"10"`
}
