package exconn

import (
	"github.com/gocql/gocql"
	"github.com/humamfauzi/go-registration/utils"
)

func ConnectToCassandra() *gocql.Session {
	cassandraEnv := utils.GetEnv("database.cassandra")
	connProfile := ComposeConnectionFromEnv(cassandraEnv, "cassandra")
	clusterAdress := strings.Split(connProfile, ",")

	cluster := gocql.NewCluster(...clusterAdress)
	cluster.Keyspace = "example"

	session, _ := cluster.CreateSession()
	return session
}
