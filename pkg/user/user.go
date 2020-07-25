package user

type User struct {
	ID                   uint64  `json:"id,omitempty"`
	Username             string  `json:"username" gorm:"unique,not null;default:null" valid:"required,alphanum"`
	PasswordHash         string  `json:"-" gorm:"not null;default: null"`
	TeamID               *uint64 `json:"team_id,omitempty"`
	Password             string  `json:"password" gorm:"-"`
	PasswordConfirmation string  `json:"password_confirmation" gorm:"-"`
}
