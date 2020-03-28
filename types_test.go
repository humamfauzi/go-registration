package registration

import (
	"testing"

	"github.com/humamfauzi/go-registration/utils"

	"github.com/humamfauzi/go-registration/exconn"
)

const ()

func TestAutoMigrateUsers(t *testing.T) {

	conn := exconn.ConnectToDB()
	defer conn.Close()
	conn.Exec("DELETE FROM users;")
	newUser := User{}

	newUser.AutoMigrate(conn)
	if !conn.HasTable(&newUser) {
		t.Error("TABLE NOT WRITTEN")
	}
}

func TestInsertUser(t *testing.T) {
	conn := exconn.ConnectToDB()
	defer conn.Close()
	token := "asdf"
	newUser := User{
		Email:    "a@a.a",
		Password: "asldkj",
		Name:     "aaa",
		Token:    &token,
	}
	dbUser := newUser.CreateUser(conn)
	if dbUser.Id == "" {
		t.Error("SHOULD HAVE A VALUE")
	}
}

func TestInsertUserBulk(t *testing.T) {
	conn := exconn.ConnectToDB()
	defer conn.Close()
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
	err := millionUser.BulkCreateUser(conn)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkInsertUserMillion(b *testing.B) {
	conn := exconn.ConnectToDB()
	defer conn.Close()
	for n := 0; n < b.N; n++ {
		millionUser := (make(Users, 100))
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
		millionUser.BulkCreateUser(conn)
	}
}
