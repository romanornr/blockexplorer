package rebuilddb

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
)

var db = database.GetDatabaseInstance()

func BuildDatabaseBlocks() {
	//database.Open()
	//database.SetupDB()

	//	blockhashChannel := make(chan []byte)
	////	blockChannel := make(chan *btcjson.GetBlockVerboseResult)
	//
	//	go func() {
	//		for i := int64(1); i < 2000; i++ {
	//			blockhashChannel <- blockdata.GetBlockHashAsync(i).CloneBytes()
	//		}
	//		close(blockhashChannel)
	//	}()
	//
	//
	//
	//	for a := range blockhashChannel {
	//		blockhash, _ := chainhash.NewHash(a)
	//		block := blockdata.GetBlock(blockhash)
	//		database.AddBlock(db, block.Hash,block )
	//	}

	////}

	progressBar := pb.StartNew(200)
	for i := int64(1); i < 200; i++ {
		//blockhash := blockdata.GetBlockHash(i) ==>
		//fmt.Println(blockdata.GetBlockHash(i))
		blockhash := blockdata.GetBlockHash(i)
		block, err := blockdata.GetBlock(blockhash)
		if err != nil {
			logs.Error(err)
		}
		database.AddBlock(db, blockhash.String(), block)
		progressBar.Increment()
		//AddIndexBlockHeightWithBlockHash(db, blockHashString, block.Height)
		//AddBlock(db, block.Hash, block) ==>
		//AddTransaction(db, block.Tx)
		//AddIndexTransactionWithBlockHash(db, blockHashString, block.Tx)

	}
	progressBar.Finish()
}
