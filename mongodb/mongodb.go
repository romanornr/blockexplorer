package mongodb

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/romanornr/cyberchain/insightjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
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
func AddBlock(Block *insightjson.BlockResult) error {
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

	err = collection.Insert(Block)

	if err != nil {
		return fmt.Errorf("Block with hash %s did not get inserted", Block.Hash)
	}

	return err
}

// get block by hash
func GetBlock(hash chainhash.Hash) (insightjson.BlockResult, error) {
	GetSession()
	collection := session.DB(Database).C("Blocks")

	result := insightjson.BlockResult{}

	err := collection.Find(bson.M{"hash": hash.String()}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

// get block by blockheight
func FetchBlockHashByBlockHeight(blockheight int64) insightjson.BlockResult {
	GetSession()
	collection := session.DB(Database).C("Blocks")
	result := insightjson.BlockResult{}

	err := collection.Find(bson.M{ "height": blockheight}).One(&result)
	if err != nil {
		panic(err)
	}

	return result
}

// get the latest block and return it in insightjson format
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

	return result, err
}

func AddTransaction(transaction *insightjson.Tx) error {
	GetSession()
	collection := session.DB(Database).C("Transactions")

	index := mgo.Index{
		Key: []string{"txid"},
		Unique: true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	err = collection.Insert(transaction)

	if err != nil {
		panic(err)
	}

	return err
}

func GetTransaction(txid chainhash.Hash) (insightjson.Tx, error) {
	GetSession()
	collection := session.DB(Database).C("Transactions")

	result := insightjson.Tx{}

	err := collection.Find(bson.M{"txid": txid.String()}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}

func AddAddressInfo(AddressInfo *insightjson.AddressInfo) error {
	GetSession()
	collection := session.DB(Database).C("AddressInfo")

	err := collection.Insert(AddressInfo)

	if err != nil {
		log.Printf("Error not able to add AddressInfo to database collection AddressInfo: %s", err)
	}

	return err
}

func GetAddressInfo(address string) (insightjson.AddressInfo, error) {
	GetSession()
	collection := session.DB(Database).C("AddressInfo")

	result := insightjson.AddressInfo{}

	err := collection.Find(bson.M{"address": address}).One(&result)
	if err != nil {
		return result, err
	}

	return result, err
}
