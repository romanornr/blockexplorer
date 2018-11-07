package mongodb

import (
	"testing"
	"log"
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
		if databases == Database {
			fmt.Println("found")
			t.Error("Old database still exists. Failed dropping.")
		}
	}
}

func TestAddBlock(t *testing.T) {

	log.Println("Adding block with height 2 to the database")

	block := insightjson.BlockResult{
		Hash: "45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
		Size: 202,
		Height: 2,
		Version: 2,
		MerkleRoot: "bf5721dfb2a84b8f39ec28dd5a7d4e8b219ca3a361dd44db2d67470818a625ea",
		Tx: []string{"bf5721dfb2a84b8f39ec28dd5a7d4e8b219ca3a361dd44db2d67470818a625ea"},
		Time: 1405608158,
		Nonce: 554156041,
		Bits: "1e01ffff",
		Difficulty: 0.001953110098725118,
		Confirmations: 257,
		PreviousBlockHash: "5ca83af67146e286610e118cc8f8e6a183c319fbb4a8fdb9e99daa2b8a29b3e3",
		NextBlockHash: "7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
		IsMainChain: true,
	}

	AddBlock(&block)

	c := session.DB(Database).C("Blocks")
	//defer session.Close()

	log.Println("Searching for block with height 2")
	result := insightjson.BlockResult{}
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

//func mockDatabase() {
//	for i := 1; i < 1001; i ++ {
//		blockhash := blockdata.GetBlockHash(int64(i))
//		block, _ := blockdata.GetBlock(blockhash)
//		newBlock,_ := insight.ConvertToInsightBlock(block)
//		AddBlock(newBlock)
//	}
//}
//
//func TestGetLastBlock(t *testing.T) {
//	mockDatabase()
//	latestblock, err := GetLastBlock()
//	if err != nil {
//		log.Printf("Error trying to find the latest block: %v", err)
//	}
//	expect := int64(1000)
//	if latestblock.Height != expect {
//		t.Errorf("Latest block is wrong \nExpected: %d \nGot: %d", expect, latestblock.Height)
//	}
//	log.Printf("Success: Latest block found with height: %d", latestblock.Height)
//}
//
//func TestAddTransaction(t *testing.T) {
//
//	hash0, _:= chainhash.NewHashFromStr("31c0cbc8411de76eac6018183e96d1cc2c904a9b50096758041eec92d9c9b9f9")
//	tx0 := blockdata.GetRawTransactionVerbose(hash0)
//	newTx0 := insight.TxConverter(tx0)
//	AddTransaction(&newTx0[0])
//
//	hash,_ := chainhash.NewHashFromStr("d78999b2ad131bd393c06738bd34996da80a556d6b1e9486447a023b91ef6ea3")
//	tx := blockdata.GetRawTransactionVerbose(hash)
//	newTx := insight.TxConverter(tx)
//	AddTransaction(&newTx[0])
//}
//
//func TestGetTransaction(t *testing.T) {
//	hash,_ := chainhash.NewHashFromStr("d78999b2ad131bd393c06738bd34996da80a556d6b1e9486447a023b91ef6ea3")
//	tx, err := GetTransaction(*hash)
//	if err != nil {
//		t.Errorf("Transaction not found with hash: %s\n", hash)
//	}
//
//	if tx.Txid != hash.String() {
//		t.Errorf("Transaction in Database got hash: %s \nExpected: %s", tx.Txid, hash.String())
//	}
//
//	log.Printf("Success: Transaction in database found with hash: %s", tx.Txid)
//}