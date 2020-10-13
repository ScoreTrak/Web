package team

import (
	"errors"
	"github.com/ScoreTrak/Web/pkg/user"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primary_key;"`

	Name string `json:"name" gorm:"unique;not null" valid:"required,alphanum"`

	Index *uint `json:"index" gorm:"unique"`

	Users []*user.User `gorm:"foreignkey:TeamID;association_foreignkey:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE" json:"-"`

	Enabled *bool `json:"enabled,omitempty" gorm:"-"`
}

func (Team) TableName() string {
	return "web_teams"
}

func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Name == "" {
		return errors.New("field Name is a mandatory parameter")
	}
	if t.ID == uuid.Nil {
		u, err := uuid.NewV4()
		if err != nil {
			return err
		}
		t.ID = u
	}
	return nil
}
