package Blockchain

import (
	"bytes"
	"encoding/gob"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"testing"
)

func BuildMockDatabase() {
	for i := int64(1); i < 10; i++ {
		blockhash := blockdata.GetBlockHash(i)
		block, _ := blockdata.GetBlock(blockhash)
		database.AddBlock(db, blockhash.String(), block)
	}
}

func TestBlockListProxy_FindBlock(t *testing.T) {

	proxy := BlockListProxy{
		Database: &BlockList{},
	}

	hashes := [5]string{"ded7508b6b6452bfc99961366e3206a7a258cf897d3148b46e590bbf6f23f3d9",
		"45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
		"7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
		"a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f",
		"d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6", //5
		//"d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88afxxxx", //invalid hash
	}

	for _, hash := range hashes {
		blockhash, err := chainhash.NewHashFromStr(hash) // check if blockhash is valid
		if err != nil {
			t.Errorf("The hash %s is not valid\n", blockhash.String())
		}
		block, err := proxy.FindBlock(blockhash)
		if err != nil {
			t.Logf("Got error: %s\n", err)
		}
		if block == nil {
			t.Errorf("Block with hash: %s is not in the database", block.Hash)
		}
	}
}

func TestBlockList_FindBlockByRPC(t *testing.T) {
	proxy := BlockListProxy{
		Database: &BlockList{},
	}

	hash := "ec02e1f752d293c2daefd4c0f66801df8cb6ee602bb1ccf219b0c55b55b123a2"
	blockhash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		t.Errorf("invalid hash")
	}

	block, err := proxy.RPC.FindBlockByRPC(blockhash)
	if err != nil {
		t.Errorf("Block %s not found but it should be found", blockhash)
	}

	if block == nil {
		t.Errorf("Block %s is empty", blockhash)
	}
}

func TestBlockList_FindBlock(t *testing.T) {
	proxy := BlockListProxy{
		Database: &BlockList{},
	}

	hash := "d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6"
	block := proxy.Database.FindBlockInDatabase(hash)
	if block == nil {
		t.Errorf("Block %s not found", hash)
	}
}

func TestAddBlockToDatabase(t *testing.T) {
	proxy := BlockListProxy{
		Database: &BlockList{},
	}

	hash := "ec02e1f752d293c2daefd4c0f66801df8cb6ee602bb1ccf219b0c55b55b123a2"
	blockhash, _ := chainhash.NewHashFromStr(hash)

	block, err := proxy.RPC.FindBlockByRPC(blockhash)
	if err != nil {
		t.Errorf("Block %s not found but it should be found", blockhash)
	}

	//proxy.Database.addBlock(block)
	proxy.AddBlockToDatabase(block)
	blockInDatabase := database.ViewBlock(hash)
	if blockInDatabase == nil {
		t.Errorf("Added block in Database but somehow not found after adding with hash: %s\n", hash)
	}
	var blockjson *btcjson.GetBlockVerboseResult
	decoder := gob.NewDecoder(bytes.NewReader(blockInDatabase))
	decoder.Decode(&blockjson)

	if blockjson.Hash != hash {
		t.Errorf("Block with hash %s added to the database but the hash does not match with what's in the database", hash)

	}
}

func TestBlockList_Cache(t *testing.T) {
	proxy := BlockListProxy{
		StackCache: BlockListCache{},
		Stacksize:  2,
	}

	hashes := [5]string{"ded7508b6b6452bfc99961366e3206a7a258cf897d3148b46e590bbf6f23f3d9", // blockheight 1
		"45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
		"7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
		"a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f",
		"d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6", // blockheight 5
	}

	for _, hash := range hashes {
		blockhash, _ := chainhash.NewHashFromStr(hash) // check if blockhash is valid
		block, _ := proxy.FindBlock(blockhash)
		proxy.addBlockToStack(block)
	}

	stack := proxy.StackCache
	if len(proxy.StackCache) > 2 {
		t.Errorf("The StackCache is bigger than the Stacksize limit! StackCache: %d Stack limit should be: %d\n", len(proxy.StackCache), proxy.Stacksize)
	}

	if stack[0].Block.Hash != "a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f" {
		t.Errorf("Wrong hash in stack. Expected: %s Actual: %s\n", hashes[3], stack[0].Block.Hash)
	}

	if stack[1].Block.Hash != "d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6" {
		t.Errorf("Wrong hash in stack. Expected: %s Actual: %s\n", hashes[4], stack[1].Block.Hash)
	}
}
