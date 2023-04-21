package libs

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"gopkg.in/mgo.v2"
)

type mongoClient struct {
	// MongoDB : DB export for Mongo
	MongoDB *mgo.Database
	*mgo.Session
}

func newMongoClient(config config.Config) *mongoClient {
	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s", config.MongoDB.Hosts...)},
		Timeout:  5 * time.Second,
		Username: config.MongoDB.User,
		Password: config.MongoDB.Password,
	}
	if config.MongoDB.SSLEnabled {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	mgoSession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	mgoSession.SetMode(mgo.Monotonic, true)
	mgoSession.SetPoolLimit(100)
	return &mongoClient{
		mgoSession.DB(config.MongoDB.Database),
		mgoSession,
	}
}
