package rebuilddb

import (
	_ "github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
	"github.com/romanornr/cyberchain/blockdata"
	"log"
)

var db = database.GetDatabaseInstance()

func BuildDatabase() {
	////database.Open()
	////database.SetupDB()
	//
	//	blockhashChannel := make(chan []byte)
	////	blockChannel := make(chan *btcjson.GetBlockVerboseResult)
	//
	//	go func() {
	//		for i := int64(1); i < 200; i++ {
	//			blockhashChannel <- blockdata.GetBlockHashAsync(i).CloneBytes()
	//		}
	//		close(blockhashChannel)
	//	}()
	//
	//
	//
	//	for a := range blockhashChannel {
	//		blockhash, _ := chainhash.NewHash(a)
	//		block, _ := blockdata.GetBlock(blockhash)
	//		fmt.Println(block)
	//		database.AddBlock(db, block.Hash,block )
	//	}


	progressBar := pb.StartNew(1000)
	for i := int64(1); i < 1000; i++ {
		//blockhash := blockdata.GetBlockHash(i) ==>
		//fmt.Println(blockdata.GetBlockHash(i))
		blockhash := blockdata.GetBlockHash(i)
		block, err := blockdata.GetBlock(blockhash)
		if err != nil {
			log.Fatal("error")
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
