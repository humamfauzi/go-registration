package exconn

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/humamfauzi/go-registration/utils"
)

const (
	LOG_QUERY = "INSERT INTO"
)

type CassandraLog struct {
	Session  *gocql.Session
	LogQuery string
}

func (cl *CassandraLog) SendLog(input []string) bool {
	return false
}

func ConnectToCassandra() *gocql.Session {
	cassandraEnv := utils.GetEnv("database.cassandra")
	connProfile := ComposeConnectionFromEnv(cassandraEnv, "cassandra")
	clusterAdress := strings.Split(connProfile, ",")

	cluster := gocql.NewCluster(clusterAdress...)
	cluster.Keyspace = "example"

	session, _ := cluster.CreateSession()
	return session
}

func InstantiateCassandraLog() CassandraLog {
	session := ConnectToCassandra()
	return CassandraLog{
		Session:  session,
		LogQuery: LOG_QUERY,
	}
}
