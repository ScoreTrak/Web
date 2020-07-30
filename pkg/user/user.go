package user

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primary_key;"`
	Username     string    `json:"username" gorm:"unique,not null;default:null" valid:"required,alphanum"`
	PasswordHash string    `json:"-" gorm:"not null;default: null"`
	TeamID       uuid.UUID `json:"team_id,omitempty" gorm:"type:uuid"`
	Password     string    `json:"password,omitempty" gorm:"-"`
	Role         string    `json:"role" gorm:"default:'blue'"`
}

func (p *User) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		u, err := uuid.NewV4()
		if err != nil {
			return err
		}
		p.ID = u
	}
	return nil
}
