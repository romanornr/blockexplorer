package mongodb

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/insight"
	"github.com/romanornr/cyberchain/insightjson"
	"fmt"
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


func Test() {
	GetSession()

	defer session.Close()

	c := session.DB("viacoin").C("Blocks")

	index := mgo.Index{
		Key: []string{"hash"},
		Unique: true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	h := blockdata.GetBlockHash(2)
	block, _ := blockdata.GetBlock(h)
	x, _ := insight.ConvertToInsightBlock(block)

	err = c.Insert(x)

	if err != nil {
		panic(err)
	}

	result := insightjson.BlockResult{}

	err = c.Find(bson.M{"hash": x.Hash}).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("result: ", result)
}
