package Blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"log"
)

type Block struct {
	Block *btcjson.GetBlockVerboseResult
}

type BlockIndex struct {
	blockHeight int64
	BlockHash string `json:"blockHash"`
}

type BlockFinder interface {
	// might use to fetch
	FindBlock(hash *chainhash.Hash) (Block, error)
	// FindBlockByRPC(hash *chainhash.Hash) (Block)
}

type BlockList []Block
type BlockListCache []Block

type BlockIndexList []BlockIndex

var db = database.GetDatabaseInstance()

type BlockListProxy struct {
	Database            *BlockList
	RPC                 *BlockList
	StackCache          BlockListCache
	Stacksize           int
	LastSearchUsedCache bool
}

// find block by looking into the database
// if the block is not in the database, check with an RPC call if it is.
// also add it in the databse if the RPC call has a result
func (b *BlockListProxy) FindBlock(hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {

	//block := b.Database.FindBlock(hash.String())
	block := b.Database.FindBlockInDatabase(hash.String())

	if block == nil {
		blockjson, _ := b.RPC.FindBlockByRPC(hash)
		b.AddBlockToDatabase(blockjson)
		return blockjson, nil
	}

	var blockjson *btcjson.GetBlockVerboseResult
	decoder := gob.NewDecoder(bytes.NewReader(block))
	decoder.Decode(&blockjson)

	return blockjson, nil
}

// find the blockhash by using the blockheight.
// Search in the database found. If nothing is found, use an RPC call.
// By using the RPC call, it means the block was never in the database. Add it into the database.
func (b *BlockListProxy) FindBlockHash(blockheight int64) BlockIndex {
	var blockIndex BlockIndex
	blockIndex.blockHeight = blockheight

	hashString := string(b.Database.FindBlockHashInDatabase(blockheight))
	hash, err := chainhash.NewHashFromStr(hashString)
	if err != nil {
		log.Printf("Could not format string to hash: %s", err)
	}

	// if the hash isn't found in the database, find it by making a RPC call
	if hashString == "" {
		hash = b.RPC.FindBlockHashByRPC(blockheight)
		block, err := b.RPC.FindBlockByRPC(hash)
		if err != nil {
			log.Printf("Blockhash %s not found: %s", blockIndex.BlockHash, err)
		}
		b.AddBlockToDatabase(block) // since it wasn't found, add the block to the database
	}
	blockIndex.BlockHash = hash.String()

	return blockIndex
}

func (b *BlockList) FindBlockByRPC(hash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	block, err := blockdata.GetBlock(hash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (b *BlockList) FindBlockHashByRPC(blockheight int64) *chainhash.Hash {
	hash := blockdata.GetBlockHash(blockheight)
	return hash
}

// find the block in the database by giving the blockhash
func (b *BlockList) FindBlockInDatabase(hash string) []byte {
	return database.ViewBlock(db, hash)
}

func (b *BlockList) FindBlockHashInDatabase(blockheight int64) []byte {
	return database.FetchBlockHashByBlockHeight(blockheight)
}

// add block to the database and index blockheight with the hash
func (b *BlockListProxy) AddBlockToDatabase(block *btcjson.GetBlockVerboseResult) {
	if block == nil {
		return // don't add anything if block is empty
	}
	database.AddBlock(db, block.Hash, block)
	database.AddIndexBlockHeightWithBlockHash(db, block.Hash, block.Height)
}

//func (b *BlockList) addBlock(block *btcjson.GetBlockVerboseResult) {
//	database.AddBlock(db, block.Hash, block)
//}

// addBlockToStack takes the user argument and adds it to the stack in place.
// if the stack is full it removes the first element on it before adding.
func (b *BlockListProxy) addBlockToStack(block *btcjson.GetBlockVerboseResult) {
	if len(b.StackCache) >= b.Stacksize {
		//b.AddBlockToDatabase(b.StackCache[0].Block)
		b.StackCache = append(b.StackCache[1:], Block{block})
	} else {
		b.StackCache.addBlockToCache(block)
	}
}

// add a new block to the end of the Block slice
func (b *BlockListCache) addBlockToCache(newBlock *btcjson.GetBlockVerboseResult) {
	*b = append(*b, Block{newBlock})
}
