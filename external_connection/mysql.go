package registration

import (
	"database/sql"
)

func connectToDB() *sql.DB {
	mysqlEnv := GetEnv("database.mysql")
	connProfile := composeConnectionFromEnv(mysqlEnv)
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
		composed += getAdditionalDbConnectionParams(parsed)
		return composed
	default:
		panic("FAILED TO PARSE DATABASE PROFILE")
	}
	return ""
}

func getAdditionalDbConnectionParams(connectionParams map[string]string) string {
	reservedKey := StringArray([]string{"username", "password", "protocol", "address", "dbname"})
	var connParams string
	for key, value := range connectionParams {
		if !reservedKey.Includes(key) {
			connParams += key + "=" + value
		}
	}
	if connParams != "" {
		connParams = "?" + connParams
	}
	return connParams
}
