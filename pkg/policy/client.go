package policy

import (
	"sync"
	"time"
)

type Client struct {
	policy *Policy
	mu     sync.RWMutex
	repo   Repo
}

func NewPolicyClient(policy *Policy, repo Repo) *Client {
	return &Client{policy: policy, mu: sync.RWMutex{}, repo: repo}
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
		rounded = rounded.Add(time.Second * 10)
		time.Sleep(time.Until(rounded))
	}
}

func (a *Client) GetPolicy() *Policy {
	a.mu.RLock()
	p := &Policy{AllowUnauthenticatedUsers: a.policy.AllowUnauthenticatedUsers, ShowPoints: a.policy.ShowPoints}
	a.mu.RUnlock()
	return p
}
