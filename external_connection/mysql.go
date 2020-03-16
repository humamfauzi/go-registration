package external_connection

import (
	"database/sql"
	"fmt"

	"github.com/humamfauzi/go-registration/utils"
)

func ConnectToDB() *sql.DB {
	mysqlEnv := utils.GetEnv("database.mysql")
	connProfile := ComposeConnectionFromEnv(mysqlEnv)
	conn, err := sql.Open("mysql", connProfile)
	if err != nil {
		panic("CANNOT CONNECT TO DATABASE")
	}
	return conn

}

func ComposeConnectionFromEnv(connection interface{}) string {
	switch connection.(type) {
	case map[string]string:
		parsed := connection.(map[string]string)
		composed := parsed["username"] + ":"
		composed += parsed["password"] + "@"
		composed += parsed["protocol"] + "("
		composed += parsed["address"] + ")"
		composed += parsed["dbname"]
		composed += GetAdditionalDbConnectionParams(parsed)
		fmt.Println(composed)
		return composed
	default:
		panic("FAILED TO PARSE DATABASE PROFILE")
	}
}

func GetAdditionalDbConnectionParams(connectionParams map[string]string) string {
	reservedKey := utils.StringArray([]string{"username", "password", "protocol", "address", "dbname"})
	var connParams string
	for key, value := range connectionParams {
		if !reservedKey.Includes(key) {
			connParams += key + "=" + value + "&"
		}
	}
	if connParams != "" {
		connParams = "?" + connParams
		connParams = connParams[0 : len(connParams)-1]
	}
	return connParams
}
