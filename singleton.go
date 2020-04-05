package registration

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/humamfauzi/go-registration/exconn"
	"github.com/jinzhu/gorm"
)

var (
	db         *gorm.DB
	logger     *gocql.Session
	documentDb context.Context
)

func InstantiateExternalConnection() {
	db = exconn.ConnectToMySQL()
	logger = exconn.ConnectToCassandra()
	documentDb = exconn.ConnectToMongo()
}
