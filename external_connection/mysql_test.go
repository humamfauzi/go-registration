package external_connection

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const mockJSON = `{
	"username": "root",
	"password": "asdf",
	"dbname": "try1",
	"address": "localhost",
	"protocol": "tcp",
	"charset": "utf8",
	"parsetime": "True",
	"loc": "Local"
}`

func TestComposeConnectionFromEnv(t *testing.T) {
	var dbProfile map[string]interface{}
	json.Unmarshal([]byte(mockJSON), &dbProfile)

	result := ComposeConnectionFromEnv(dbProfile)
	if result == "" {
		t.Errorf("SHOULD HAVE A VALUE")
	}
}

func TestDBConnection(t *testing.T) {
	ConnectToDB()
}
