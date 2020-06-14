package config

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"os"
)

type Config struct {
	Port                 string
	MongoDBRawConnString string
	MongoDBConnString    connstring.ConnString
}

const (
	portEnvVar              = "PORT"
	defaultPort             = "8080"
	mongoDBConnStringEnvVar = "MONGODB_URI"
	defaultMongoDBURI       = "mongodb://localhost:27017/test"
)

func ReadConfig(logf func(string, ...interface{})) (Config, error) {
	var (
		c   Config
		ok  bool
		err error
	)
	if c.Port, ok = os.LookupEnv(portEnvVar); !ok {
		logf("%s not set, applying default %s", portEnvVar, defaultPort)
		c.Port = "8080"
	}
	if c.MongoDBRawConnString, ok = os.LookupEnv(mongoDBConnStringEnvVar); !ok {
		logf("%s not set, applying default %s", mongoDBConnStringEnvVar, defaultMongoDBURI)
		c.MongoDBRawConnString = defaultMongoDBURI
	}
	c.MongoDBRawConnString += "?retryWrites=false"
	if c.MongoDBConnString, err = connstring.Parse(c.MongoDBRawConnString); err != nil {
		return c, errors.Wrap(err, "failed to parse MongoDB conn string")
	}
	return c, err
}
