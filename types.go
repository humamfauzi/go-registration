package registration

import (
	"errors"
	"time"

	"github.com/humamfauzi/go-registration/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	module = "USERS"
)

type User struct {
	Id        string  `gorm:"type:varchar(100);unique_index;primary_key" json:"id"`
	Email     string  `gorm:"type:varchar(255)"`
	Password  string  `gorm:"type:varchar(255)"`
	Name      string  `gorm:"type:varchar(255)`
	Token     *string `gorm:"type:varchar(255)`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) SetToken(token string) {
	u.Token = &token
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) SetId(id string) {
	u.Id = id
}

func (u User) AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&u)
}

func (u User) CreateUser(db *gorm.DB) User {
	uniqId := utils.GenerateUUID("user", 2)
	u.SetId(uniqId)
	db.Create(&u)
	return u
}

func (u User) DeleteUser(db *gorm.DB) (bool, error) {
	if !u.hasID() {
		return false, errors.New("Deletion should have an id")
	}
	db.Delete(&u)
	return true, nil
}

func (u User) hasID() bool {
	if u.Id == "" {
		return false
	} else {
		return true
	}
}

func (u User) UpdateUser(db *gorm.DB, newUserProfile User) error {
	if newUserProfile.hasID() {
		return errors.New("cannot update with id")
	}
	db.Model(&u).Update(newUserProfile)
	return nil
}

type Users []User

func (u Users) hasSomeUserID() bool {
	for i := 0; i < len(u); i++ {
		if u[i].Id != "" {
			return true
		}
	}
	return false
}

func (u Users) BulkUpdateUser(db *gorm.DB, newUsersProfile Users) error {
	if u.hasSomeUserID() {
		return errors.New("updated profile should not have an id")
	}
	for index := 0; index < len(u); index++ {
		db.Model(&u).Update(u[index])
	}
	return nil
}

func (u Users) BulkCreateUser(db *gorm.DB) error {
	for index := 0; index < len(u); index++ {
		uniqId := utils.GenerateUUID("user", 2)
		u[index].SetId(uniqId)
		db.Create(&u[index])
	}
	return nil
}

func generateEmptyUser() User {
	var user User
	return user
}
