package mongodb

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/globalsign/mgo"
	"github.com/romanornr/cyberchain/insight"
)

const (
	MongoDBHosts = "localhost"
	Database     = "Viacoin"
)

var session *mgo.Session

var mongoDBDialInfo = &mgo.DialInfo{
	Addrs:    []string{MongoDBHosts},
	Timeout: 60 * time.Second,
	Database: Database,
}

func GetSession() *mgo.Session {
	if session == nil {
		var err error
		session, err = mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			session.Close()
			panic(err)
		}
	}
	return session
}

// delete the database. Only use for testing
func DropDatabase() error {
	err := session.DB(Database).DropDatabase()
	if err != nil {
		panic(err)
	}
	return err
}

// add Blocks to the database. Collection name: Blocks
// NEED FIX: DOES NOT ERROR WHEN INSERTING EXISTING KEY :S
func AddBlock(Block *btcjson.GetBlockVerboseResult) error {
	GetSession()
	//defer session.Close()

	collection := session.DB("viacoin").C("Blocks")

	index := mgo.Index{
		Key:    []string{"hash"},
		Unique: true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = collection.Insert()

	insightBlock, _ := insight.ConvertToInsightBlock(Block)

	err = collection.Insert(insightBlock)

	if err != nil {
		return fmt.Errorf("Block with hash %s did not get inserted", Block.Hash)
	}

	return err
}
