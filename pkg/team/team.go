package team

import (
	"github.com/ScoreTrak/Web/pkg/user"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID uuid.UUID `json:"id,omitempty" gorm:"type:uuid;primary_key;"`

	Name string `json:"name" gorm:"unique;not null;default:null" valid:"required,alphanum"`

	Users []*user.User `gorm:"foreignkey:TeamID;association_foreignkey:ID" json:"-"`

	Enabled *bool `json:"enabled,omitempty" gorm:"-"`
}

func (p *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		u, err := uuid.NewV4()
		if err != nil {
			return err
		}
		p.ID = u
	}
	return nil
}
