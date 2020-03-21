package external_connection

import (
	"database/sql"

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
	case map[string]interface{}:
		parsed := connection.(map[string]interface{})
		composed := utils.InterpretInterfaceString(parsed["username"], "root")
		composed += utils.InterpretInterfaceString(parsed["password"], "")
		composed += utils.InterpretInterfaceString(parsed["protocol"], "tcp")
		composed += utils.InterpretInterfaceString(parsed["adress"], "localhost")
		composed += utils.InterpretInterfaceString(parsed["dbname"], "try1")
		composed += GetAdditionalDbConnectionParams(parsed)
		return composed
	default:
		panic("FAILED TO PARSE DATABASE PROFILE")
	}
}

func GetAdditionalDbConnectionParams(connectionParams map[string]interface{}) string {
	reservedKey := utils.StringArray([]string{"username", "password", "protocol", "address", "dbname"})
	var connParams string
	for key, value := range connectionParams {
		if !reservedKey.Includes(key) {
			connParams += key + "=" + value.(string) + "&"
		}
	}
	if connParams != "" {
		connParams = "?" + connParams
		connParams = connParams[0 : len(connParams)-1]
	}
	return connParams
}
