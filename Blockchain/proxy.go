package Blockchain

import (
	"github.com/btcsuite/btcd/btcjson"
	"fmt"
	"github.com/romanornr/cyberchain/database"
	"encoding/gob"
	"bytes"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Block struct {
	Block *btcjson.GetBlockVerboseResult
}

type BlockFinder interface {
	FindBlock(hash string)
}

type BlockList []Block

var db = database.GetDatabaseInstance()

type BlockListProxy struct {
	Database *BlockList
}

// find block by looking into the database
// if the block is not in the database, check with an RPC call if it is.
// also add it in the databse if the RPC call has a result
func (b *BlockListProxy) FindBlock(hash string) (*btcjson.GetBlockVerboseResult, error) {

	block := b.Database.FindBlock(hash)


	if block == nil {
		fmt.Println("Checking with RPC call now")
		chainhash, err := chainhash.NewHashFromStr(hash)
		if err != nil {
			return Block{}.Block, fmt.Errorf("The hash %s is not valid", chainhash)
		}
		blockjson := blockdata.GetBlock(chainhash)
		b.AddBlockToDatabase(blockjson)
		return blockjson, nil
	}

	var blockjson *btcjson.GetBlockVerboseResult
	decoder := gob.NewDecoder(bytes.NewReader(block))
	decoder.Decode(&blockjson)

	return blockjson, nil

}

// find the block in the database by giving the blockhash
func (l *BlockList) FindBlock(hash string) []byte {
	return database.ViewBlock(hash)
}

func (b *BlockListProxy) AddBlockToDatabase(block *btcjson.GetBlockVerboseResult) {
	b.Database.addBlock(block)
	//can do something like b.database.addTransaction(block)
}

func (b *BlockList) addBlock(block *btcjson.GetBlockVerboseResult) {
	database.AddBlock(db, block.Hash, block)
}