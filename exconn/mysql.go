package exconn

import (
	"github.com/humamfauzi/go-registration/utils"
	"github.com/jinzhu/gorm"
)

func ConnectToMySQL() *gorm.DB {
	mysqlEnv := utils.GetEnv("database.mysql")
	connProfile := ComposeConnectionFromEnv(mysqlEnv, "mysql")
	conn, err := gorm.Open("mysql", connProfile)
	if err != nil {
		panic(err)
	}
	return conn

}
func CloseConnectionMySQL(db *gorm.DB) {
	db.Close()
}

func GetAdditionalDbConnectionParams(connectionParams map[string]interface{}) string {
	reservedKey := utils.StringArray([]string{"username", "password", "protocol", "address", "dbname"})
	var connParams string
	for key, value := range connectionParams {
		if !reservedKey.Includes(key) {
			if key == "parsetime" {
				key = "parseTime"
			}
			connParams += key + "=" + value.(string) + "&"
		}
	}
	if connParams != "" {
		connParams = "?" + connParams
		connParams = connParams[0 : len(connParams)-1]
	}
	return connParams
}
