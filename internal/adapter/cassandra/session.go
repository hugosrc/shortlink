package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/hugosrc/shortlink/internal/util"
	"github.com/spf13/viper"
)

func New(conf viper.Viper) (*gocql.Session, error) {
	cluster := gocql.NewCluster(conf.GetString("CASSANDRA_SERVER"))
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.GetString("CASSANDRA_USER"),
		Password: conf.GetString("CASSANDRA_PASSWORD"),
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, util.WrapErrorf(err, util.ErrCodeUnknown, "error connecting to cassandra server")
	}

	return session, nil
}
