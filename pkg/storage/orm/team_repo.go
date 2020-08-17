package orm

import (
	"errors"
	"fmt"
	"github.com/ScoreTrak/ScoreTrak/pkg/logger"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type teamRepo struct {
	db  *gorm.DB
	log logger.LogInfoFormat
}

func NewTeamRepo(db *gorm.DB, log logger.LogInfoFormat) team.Repo {
	return &teamRepo{db, log}
}

func (h *teamRepo) Delete(id uuid.UUID) error {
	h.log.Debugf("deleting the team with id : %h", id)
	result := h.db.Delete(&team.Team{}, "id = ?", id)
	if result.Error != nil {
		errMsg := fmt.Sprintf("error while deleting the team with id : %d", id)
		h.log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (h *teamRepo) GetAll() ([]*team.Team, error) {
	h.log.Debug("get all the teams")
	teams := make([]*team.Team, 0)
	err := h.db.Find(&teams).Error
	if err != nil {
		h.log.Debug("not a single team found")
		return nil, err
	}
	return teams, nil
}

func (h *teamRepo) GetByID(id uuid.UUID) (*team.Team, error) {
	h.log.Debugf("get team details by id : %h", id)
	tm := &team.Team{}
	err := h.db.Where("id = ?", id).First(tm).Error
	if err != nil {
		h.log.Errorf("team not found with id : %h, reason : %v", id, err)
		return nil, err
	}
	return tm, nil
}

func (h *teamRepo) Store(tm []*team.Team) error {
	err := h.db.Create(tm).Error
	if err != nil {
		h.log.Errorf("error while creating the team, reason : %v", err)
		return err
	}
	return nil
}

func (h *teamRepo) Update(tm *team.Team) error {
	h.log.Debugf("updating the team, id : %v", tm.ID)
	err := h.db.Model(tm).Updates(team.Team{Name: tm.Name, Enabled: tm.Enabled, Index: tm.Index}).Error //TODO: Adjust Casbin rules on TeamID, change
	if err != nil {
		h.log.Errorf("error while updating the team, reason : %v", err)
		return err
	}
	return nil
}
