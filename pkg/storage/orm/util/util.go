package util

import (
	"github.com/ScoreTrak/Web/pkg/image"
	"github.com/ScoreTrak/Web/pkg/policy"
	"github.com/ScoreTrak/Web/pkg/role"
	"github.com/ScoreTrak/Web/pkg/team"
	"github.com/ScoreTrak/Web/pkg/user"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var uuid1 = uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001")

func CreateAllTables(db *gorm.DB) (err error) {
	err = db.AutoMigrate(&policy.Policy{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&team.Team{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&image.Image{})
	if err != nil {
		return err
	}
	return nil
}

func CreateBlackTeam(db *gorm.DB) (err error) {
	err = db.Create([]*team.Team{{ID: uuid1, Name: "Black Team"}}).Error
	if err != nil {
		serr, ok := err.(*pgconn.PgError)
		if !ok || serr.Code != "23505" {
			return err
		}
	}
	return nil
}

func CreateAdminUser(db *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = db.Create([]*user.User{{ID: uuid1, TeamID: uuid1, Username: "admin", Role: role.Black, PasswordHash: string(hashedPassword)}}).Error
	if err != nil {
		serr, ok := err.(*pgconn.PgError)
		if !ok || serr.Code != "23505" {
			return err
		}
	}
	return nil
}

func CreatePolicy(db *gorm.DB) (*policy.Policy, error) {
	p := &policy.Policy{ID: 1}
	err := db.Create(p).Error
	if err != nil {
		serr, ok := err.(*pgconn.PgError)
		if !ok {
			if serr.Code != "23505" {
				panic(err)
			} else {
				db.Take(p)
			}
		}
	}
	return p, nil
}
