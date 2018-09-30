package rebuilddb

import (
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
	"github.com/romanornr/cyberchain/blockdata"
	"log"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"sync"
)

var db = database.GetDatabaseInstance()

/*
	note: 2000 blocks costs currently 8.4 MB and ~39 seconds to save. Running into performance issues.
	goroutine "go addTransactions(block) speed up from ~39 seconds to ~29 seconds. 25% speed up
 */
func BuildDatabase() {
	progressBar := pb.StartNew(2000)
	wg := sync.WaitGroup{}
	wg.Add(1999)

	for i := int64(1); i < 2000; i++ {
		go resolveBlockToDB(i, progressBar, &wg)
	}

	wg.Wait()
	progressBar.Finish()
}

func resolveBlockToDB(i int64, prBar *pb.ProgressBar, callerWG *sync.WaitGroup) {
	blockhash := blockdata.GetBlockHash(i)
	block, err := blockdata.GetBlock(blockhash)
	if err != nil {
		log.Fatal(err)
		return
	}

	localWG := sync.WaitGroup{}
	localWG.Add(3)

	// adding transactions to the "Transactions" bucket
	go func() {
		defer localWG.Done()
		for j := 0; j < len(block.Tx); j++ {
			txhash, _ := chainhash.NewHashFromStr(block.Tx[j])
			tx := blockdata.GetRawTransactionVerbose(txhash)
			database.AddTransaction(db, tx)
		}
	}()

	go func() {
		defer localWG.Done()

		// bucket:"Blocks"  key:blockhash  value:blockVerboseResult
		database.AddBlock(db, blockhash.String(), block)
	}()

	go func() {
		defer localWG.Done()

		// bucket:"Blockheight" key:blockheight value:blockhash
		database.AddIndexBlockHeightWithBlockHash(db, block.Hash, block.Height)
	}()

	localWG.Wait()
	callerWG.Done()
	prBar.Increment()
}

// future opt might be using hex instead of TransactionVerbose
// 2000 blocks takes the same size but 35 seconds instead of 39
// for j := 0; j < len(block.Tx); j++ {
// 	txhash, _ := chainhash.NewHashFromStr(block.Tx[j])
// 	r := blockdata.GetRawTransaction(txhash)
// 	txHex, _ := hex.DecodeString("010000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff3503d3fa520fe4b883e5bda9e7a59ee4bb99e9b1bc205b323031382d30392d31355431343a32343a32392e3237343736363538315a5dffffffff0294357700000000001976a914b722782e401ed7a31135580ada74962ef32e5cd288ac0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000")
//
// 	database.AddRawTransaction(db, txHex, r)
// }
