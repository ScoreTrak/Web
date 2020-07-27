package orm

import (
	"errors"
	"fmt"
	"github.com/L1ghtman2k/ScoreTrak/pkg/logger"
	"github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"
	"gorm.io/gorm"
)

type userRepo struct {
	db  *gorm.DB
	log logger.LogInfoFormat
}

func NewUserRepo(db *gorm.DB, log logger.LogInfoFormat) user.Repo {
	return &userRepo{db, log}
}

func (h *userRepo) Delete(id uint64) error {
	h.log.Debugf("deleting the user with id : %h", id)
	result := h.db.Delete(&user.User{}, "id = ?", id)
	if result.Error != nil {
		errMsg := fmt.Sprintf("error while deleting the user with id : %d", id)
		h.log.Errorf(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (h *userRepo) GetAll() ([]*user.User, error) {
	h.log.Debug("get all the users")
	users := make([]*user.User, 0)
	err := h.db.Find(&users).Error
	if err != nil {
		h.log.Debug("not a single user found")
		return nil, err
	}
	return users, nil
}

func (h *userRepo) GetByID(id uint64) (*user.User, error) {
	h.log.Debugf("get user details by id : %h", id)
	usr := &user.User{}
	err := h.db.Where("id = ?", id).First(usr).Error
	if err != nil {
		h.log.Errorf("user not found with id : %h, reason : %v", id, err)
		return nil, err
	}
	return usr, nil
}

func (h *userRepo) GetByUsername(username string) (*user.User, error) {
	h.log.Debugf("get user details by id : %h", username)
	usr := &user.User{}
	err := h.db.Where("username = ?", username).First(usr).Error
	if err != nil {
		h.log.Errorf("user not found with id : %h, reason : %v", username, err)
		return nil, err
	}
	return usr, nil
}

func (h *userRepo) Store(usr *user.User) error {
	h.log.Debugf("creating the user with id : %v", usr.ID)
	err := h.db.Create(usr).Error
	if err != nil {
		h.log.Errorf("error while creating the user, reason : %v", err)
		return err
	}
	return nil
}

func (h *userRepo) Update(usr *user.User) error {
	h.log.Debugf("updating the user, id : %v", usr.ID)
	err := h.db.Model(usr).Updates(user.User{PasswordHash: usr.PasswordHash, Username: usr.Username, TeamID: usr.TeamID}).Error //TODO: Adjust Casbin rules on TeamID, change
	if err != nil {
		h.log.Errorf("error while updating the user, reason : %v", err)
		return err
	}
	return nil
}
