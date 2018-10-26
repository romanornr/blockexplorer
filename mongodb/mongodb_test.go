package mongodb

import (
	"testing"
	"log"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/insightjson"
	"github.com/globalsign/mgo/bson"
	"fmt"
)

func TestGetSession(t *testing.T) {
	log.Println("Getting session to mongodb")
	GetSession()
}

// drop if database exist so tests can start clean
func TestDropDatabase(t *testing.T) {
	session := GetSession()

	log.Println("Dropping old existing database")
	DropDatabase()

	databases, _ := session.DatabaseNames()

	for _, databases := range databases {
		if databases == "minimalGraphql" {
			fmt.Println("found")
		}
	}
}

func TestAddBlock(t *testing.T) {

	log.Println("Get block with height 2 via an RPC call")
	hash := blockdata.GetBlockHash(2)
	block, _ := blockdata.GetBlock(hash)

	log.Println("Adding block to the database")
	AddBlock(block)

	result := insightjson.BlockResult{}

	c := session.DB("viacoin").C("Blocks")
	defer session.Close()

	log.Println("Searching for block with height 2")
	err := c.Find(bson.M{"hash": block.Hash}).One(&result)
	if err != nil {
		panic(err)
	}
	expect := "45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc"

	if result.Hash != expect {
		t.Errorf("Expected: %s \nGot: %s\n", expect, result.Hash)
	}

	log.Printf("Success: Block with hash %s found in the database", result.Hash)
}
