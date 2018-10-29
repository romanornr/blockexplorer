package mongodb

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/globalsign/mgo"
	"github.com/romanornr/cyberchain/insight"
	"github.com/globalsign/mgo/bson"
	"github.com/romanornr/cyberchain/insightjson"
)

const (
	MongoDBHosts = "localhost"
	Database     = "viacoin"
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

	collection := session.DB(Database).C("Blocks")

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

// get block by blockheight
func FetchBlockHashByBlockHeight(blockheight int64) insightjson.BlockResult {
	GetSession()
	collection := session.DB(Database).C("Blocks")
	result := insightjson.BlockResult{}
	//err := collection.Find( { "Height" : blockheight})
	err := collection.Find(bson.M{ "height": blockheight}).One(&result)
	if err != nil {
		panic(err)
	}

	return result
}

func GetLastBlock() (insightjson.BlockResult, error) {
	GetSession()
	collection := session.DB(Database).C("Blocks")
	result := insightjson.BlockResult{}

	dbSize, err := collection.Count()
	if err != nil {
		return result, err
	}


	err = collection.Find(nil).Skip(dbSize-1).One(&result)
	if err != nil {
		return result, err
	}

	fmt.Println("last block:")

	fmt.Println(result)

	return result, err
}

