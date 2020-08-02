package orm

import (
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/pkg/policy"
	"gorm.io/gorm"
)

type policyRepo struct {
	db  *gorm.DB
	log logger.LogInfoFormat
}

func NewPolicyRepo(db *gorm.DB, log logger.LogInfoFormat) policy.Repo {
	return &policyRepo{db, log}
}

func (h *policyRepo) Get() (*policy.Policy, error) {
	p := &policy.Policy{}
	p.ID = 1
	h.db.Take(p)
	return p, nil
}

func (h *policyRepo) Update(tm *policy.Policy) error {
	h.log.Debugf("updating the policy")
	tm.ID = 1
	err := h.db.Model(tm).Updates(policy.Policy{AllowUnauthenticatedUsers: tm.AllowUnauthenticatedUsers, ShowPoints: tm.ShowPoints}).Error
	if err != nil {
		h.log.Errorf("error while updating the config, reason : %v", err)
		return err
	}
	return nil
}
