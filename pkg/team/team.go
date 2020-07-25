package team

import "github.com/L1ghtman2k/ScoreTrakWeb/pkg/user"

type Team struct {
	ID uint64 `json:"id" gorm:"primary_key"`

	Name string `json:"name" gorm:"unique;not null;default:null" valid:"required,alphanum"`

	Users []*user.User `gorm:"foreignkey:TeamID;association_foreignkey:ID" json:"-"`

	Enabled *bool `json:"enabled,omitempty" gorm:"-"`
}
