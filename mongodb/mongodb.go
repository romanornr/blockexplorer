package mongodb

import (
	"github.com/globalsign/mgo"
)

var session *mgo.Session

func GetSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.Dial("mongodb://localhost")
		if err != nil {
			session.Close()
			panic(err)
		}
	}
	return session
}
