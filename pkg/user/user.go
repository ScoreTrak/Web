package user

import (
	"errors"
	"github.com/ScoreTrak/Web/pkg/role"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primary_key;"`
	Username     string    `json:"username" gorm:"unique,not null;default:null" valid:"required,alphanum"`
	PasswordHash string    `json:"password_hash" gorm:"not null;default: null"`
	TeamID       uuid.UUID `json:"team_id,omitempty" gorm:"type:uuid"`
	Password     string    `json:"password,omitempty" gorm:"-"`
	Role         string    `json:"role" gorm:"default:'blue'"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		uid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		u.ID = uid
	}
	return nil
}

func (u User) Validate(db *gorm.DB) {
	if u.Role != "" && u.Role != role.Black && u.Role != role.Blue {
		db.AddError(errors.New("you must specify a correct role"))
	}
}
