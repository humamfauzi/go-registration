package exconn

import (
	"github.com/humamfauzi/go-registration/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectToPostgres() *gorm.DB {
	postgresEnv := utils.GetEnv("database.postgres")
	connProfile := ComposeConnectionFromEnv(postgresEnv, "postgres")
	if err != nil {
		panic(err)
	}
	return conn
}

func CloseConnectionPostgres(db *gorm.DB) {
	db.Close()
}
