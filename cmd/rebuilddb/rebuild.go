package rebuilddb

import (
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
)

func BuildDatabaseBlocks() {
	database.Open()
	database.SetupDB()
	for i := int64(1); i < 100000; i++ {
		//blockhash := blockdata.GetBlockHash(i) ==>
		//fmt.Println(blockdata.GetBlockHash(i))
		blockdata.GetBlockHash(i)
		//block := blockdata.GetBlock(blockhash) ==>
		//AddIndexBlockHeightWithBlockHash(db, blockHashString, block.Height)
		//AddBlock(db, block.Hash, block) ==>
		//AddTransaction(db, block.Tx)
		//AddIndexTransactionWithBlockHash(db, blockHashString, block.Tx)
	}
}

