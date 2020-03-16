package external_connection

import (
	"encoding/json"
	"testing"
)

const mockJSON = `{
	"username": "root",
	"password": "asdf",
	"database": "try1",
	"address": "localhost",
	"protocol": "tcp",
	"charset": "utf8",
	"parsetime": "True",
	"loc": "Local"
}`

func TestComposeConnectionFromEnv(t *testing.T) {
	var dbProfile map[string]string
	json.Unmarshal([]byte(mockJSON), &dbProfile)

	result := ComposeConnectionFromEnv(dbProfile)
	if result == "" {
		t.Errorf("SHOULD HAVE A VALUE")
	}

	expectedResult := "root:asdf@tcp(localhost)?database=try1&charset=utf8&parsetime=True&loc=Local"
	if result != expectedResult {
		t.Errorf("Expected: %s, Got: %s", result, expectedResult)
	}
}

// func TestConnectToDB(t *testing.T) {

// }
