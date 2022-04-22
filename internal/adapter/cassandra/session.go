package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/spf13/viper"
)

func New(conf viper.Viper) (*gocql.Session, error) {
	cluster := gocql.NewCluster(conf.GetString("CASSANDRA_SERVER"))
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.GetString("CASSANDRA_USER"),
		Password: conf.GetString("CASSANDRA_PASSWORD"),
	}

	return cluster.CreateSession()
}
