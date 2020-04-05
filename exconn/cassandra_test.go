package exconn

import (
	"fmt"
	"os"
	"testing"

	"github.com/gocql/gocql"
)

var (
	session *gocql.Session
)

func TestMain(m *testing.M) {
	Setup()
	code := m.Run()
	Teardown()
	os.Exit(code)
}

func Setup() {
	session = ConnectToCassandra()

	qu := session.Query(`SELECT table_name FROM system_schema.tables`)
	iter := qu.Iter()
	var result map[string]interface{}
	ok := iter.MapScan(result)
	if !ok {
		panic("PANIC")
	} else {
		fmt.Println(result)
	}
	createUsersTableQuery := `CREATE TABLE users(
		user_id ascii PRIMARYKEY,
		email ascii,
		password ascii,
		phone ascii,
		token ascii
	)
	`
	err := session.Query(createUsersTableQuery).Exec()
	if err != nil {
		panic(err)
	}
	insertUserQuery := `INSERT INTO users (user_id, email, phone) VALUES (user, pass@pass.com, 09182301293)`
	err = session.Query(insertUserQuery).Exec()
	if err != nil {
		panic(err)
	}

}

func Teardown() {
	session.Close()
}
func TestConnectToCassandra(t *testing.T) {
	ConnectToCassandra()
}
