package rebuilddb

import (
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
	"log"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

var db = database.GetDatabaseInstance()

func BuildDatabase() {

	progressBar := pb.StartNew(2000)
	for i := int64(1); i < 2000; i++ {

		blockhash := blockdata.GetBlockHash(i)
		block, err := blockdata.GetBlock(blockhash)
		if err != nil {
			log.Fatal("error")
		}

		database.AddBlock(db, blockhash.String(), block)                        // bucket:"Blocks"  key:blockhash  value:blockVerboseResult
		database.AddIndexBlockHeightWithBlockHash(db, block.Hash, block.Height) //bucket:"Blockheight" key:blockheight value:blockhash

		//adding transactions the "Transactions" bucket
		for j := 0; j < len(block.Tx); j++ {
			txhash, _ := chainhash.NewHashFromStr(block.Tx[j])
			tx := blockdata.GetRawTransactionVerbose(txhash)
			database.AddTransaction(db, tx)
		}

		// note: 2000 blocks costs currently 8.4 MB and ~39 seconds to save. Running into performance issues.
		progressBar.Increment()
	}
	progressBar.Finish()
}


// future opt might be using hex instead of TransactionVerbose
// 2000 blocks takes the same size but 35 seconds instead of 39
//for j := 0; j < len(block.Tx); j++ {
//	txhash, _ := chainhash.NewHashFromStr(block.Tx[j])
//	r := blockdata.GetRawTransaction(txhash)
//	txHex, _ := hex.DecodeString("010000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff3503d3fa520fe4b883e5bda9e7a59ee4bb99e9b1bc205b323031382d30392d31355431343a32343a32392e3237343736363538315a5dffffffff0294357700000000001976a914b722782e401ed7a31135580ada74962ef32e5cd288ac0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000")
//
//	database.AddRawTransaction(db, txHex, r)
//}
