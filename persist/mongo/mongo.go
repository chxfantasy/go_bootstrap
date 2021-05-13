package mongo

import (
	"errors"

	mgo "github.com/globalsign/mgo"
)

// ConfigDef mongo ConfigDef
type ConfigDef struct {
	Dsn string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
}

// NewMongoSession NewMongoSession
func NewMongoSession(mgConf *ConfigDef) (mongoSession *mgo.Session, err error) {
	if mgConf == nil {
		return nil, errors.New("nil mongo.ConfigDef")
	}
	if mgConf.Dsn == "" {
		return nil, errors.New("nil mongo dsn")
	}
	mongoSession, err = mgo.Dial(mgConf.Dsn)
	if err != nil {
		return
	}
	if mongoSession == nil {
		return nil, errors.New("get nil mongoSession")
	}
	//mongoSession.SetMode(mgo.SecondaryPreferred, true)
	return
}
