package policy

type Policy struct {
	ID                                 uint  `json:"-" gorm:"primary_key;"`
	AllowUnauthenticatedUsers          *bool `json:"allow_unauthenticated_users" gorm:"not null;default: false"`
	AllowChangingUsernamesAndPasswords *bool `json:"allow_changing_usernames_and_passwords" gorm:"not null;default: false"`
	ShowPoints                         *bool `json:"allow_to_see_points" gorm:"not null;default: false"`
}
