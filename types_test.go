package registration

import (
	"os"
	"testing"

	"github.com/humamfauzi/go-registration/utils"
	"github.com/jinzhu/gorm"
)

var (
	conn *gorm.DB
)

func TestMain(m *testing.M) {
	Setup()
	code := m.Run()
	Teardown()
	os.Exit(code)
}

func Setup() {
	InstantiateExternalConnection()
}

func Teardown() {
	conn.Close()
}

func TestAutoMigrateUsers(t *testing.T) {
	conn.Exec("DELETE FROM users;")
	newUser := User{}

	newUser.AutoMigrate()
	if !conn.HasTable(&newUser) {
		t.Error("TABLE NOT WRITTEN")
	}
}

func TestInsertUser(t *testing.T) {
	token := "asdf"
	newUser := User{
		Email:    "a@a.a",
		Password: "asldkj",
		Name:     "aaa",
		Token:    &token,
	}
	dbUser := newUser.CreateUser()
	if dbUser.Id == "" {
		t.Error("SHOULD HAVE A VALUE")
	}
}

func TestInsertUserBulk(t *testing.T) {
	millionUser := (make(Users, 10))
	var user User
	for i := 0; i < 10; i++ {
		var token string
		user = User{
			Email:    utils.GenerateRandomString(5) + "@" + utils.GenerateRandomString(5) + ".com",
			Password: utils.GenerateRandomString(12),
			Name:     utils.GenerateRandomString(4),
			Token:    &token,
		}
		millionUser[i] = user
	}
	err := millionUser.BulkCreateUser()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkInsertUserMillion(b *testing.B) {
	for n := 0; n < b.N; n++ {
		millionUser := (make(Users, 10000))
		var user User
		for i := 0; i < 100; i++ {
			user = User{
				Email:    utils.GenerateRandomString(5) + "@" + utils.GenerateRandomString(5) + ".com",
				Password: utils.GenerateRandomString(12),
				Name:     utils.GenerateRandomString(4),
				Token:    nil,
			}
			millionUser[i] = user
		}
		millionUser.BulkCreateUser()
	}
}

func BenchmarkAppendString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var someString string
		for i := 0; i < 100; i++ {
			someString += "a"
		}
	}
}

func BenchmarkArrayString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		someString := make([]rune, 100)
		for i := 0; i < 100; i++ {
			someString[i] = 'a'
		}
		_ = string(someString)

	}
}
