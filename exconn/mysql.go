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
