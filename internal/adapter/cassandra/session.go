package cassandra

import "github.com/gocql/gocql"

const (
	cassandraServer   string = "127.0.0.1:9042"
	cassandraUser     string = "cassandra"
	cassandraPassword string = "cassandra"
)

func New() (*gocql.Session, error) {
	cluster := gocql.NewCluster(cassandraServer)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cassandraUser,
		Password: cassandraPassword,
	}

	return cluster.CreateSession()
}
