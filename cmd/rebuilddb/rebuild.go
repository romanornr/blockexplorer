package rebuilddb

import (
	_ "github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
	"log"
)

var db = database.GetDatabaseInstance()

func BuildDatabase() {

	progressBar := pb.StartNew(1000)
	for i := int64(1); i < 1000; i++ {

		blockhash := blockdata.GetBlockHash(i)
		block, err := blockdata.GetBlock(blockhash)
		if err != nil {
			log.Fatal("error")
		}

		database.AddBlock(db, blockhash.String(), block)                        // bucket:"Blocks"  key:blockhash  value:blockVerboseResult
		database.AddIndexBlockHeightWithBlockHash(db, block.Hash, block.Height) //bucket:"Blockheight" key:blockheight value:blockhash

		progressBar.Increment()
	}

	progressBar.Finish()

	// example txhex to tx
	//a, _ := hex.DecodeString("010000000001010000000000000000000000000000000000000000000000000000000000000000ffffffff3503d3fa520fe4b883e5bda9e7a59ee4bb99e9b1bc205b323031382d30392d31355431343a32343a32392e3237343736363538315a5dffffffff0294357700000000001976a914b722782e401ed7a31135580ada74962ef32e5cd288ac0000000000000000266a24aa21a9ede2f61c3f71d1defd3fa999dfa36953755c690689799962b48bebd836974e8cf90120000000000000000000000000000000000000000000000000000000000000000000000000")
	//b, _ := btcutil.NewTxFromBytes(a)
	//k := b.MsgTx().TxIn[0].SerializeSize()
	//fmt.Println(string(k))

}
