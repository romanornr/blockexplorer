package rebuilddb

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/address"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
	"log"
	"sync"
	"encoding/gob"
	"bytes"
	"fmt"
)

var db = database.GetDatabaseInstance()

/*
	note: 2000 blocks costs currently 8.4 MB and ~39 seconds to save. Running into performance issues.
	goroutine "go addTransactions(block) speed up from ~39 seconds to ~29 seconds. 25% speed up

	try to analyze this address: https://chainz.cryptoid.info/via/address.dws?369935.htm
*/


func BuildDatabase() {

	currentBlock := 5473972
	progressBar := pb.StartNew(currentBlock)
	wg := sync.WaitGroup{}
	wg.Add(currentBlock-1)

	for i := int64(5422071); i < int64(currentBlock); i++ {
		go resolveBlockToDB(i, progressBar, &wg)
	}

	wg.Wait()
	progressBar.Finish()
	test()
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

// index adresses by calculating every transaction, balance etc
// example: https://explorer.viacoin.org/api/addr/Vrh9ro5WhykxrPPBe2cgyNiB2sAVqzkWjX
func resolveAddresses(transaction *btcjson.TxRawResult) {

	var address address.Index

	for i := 0; i < len(transaction.Vout); i++ {
		for j := 0; i < len(transaction.Vout[i].ScriptPubKey.Addresses); j++ {
			address.AddrStr = transaction.Vout[i].ScriptPubKey.Addresses[j]
			address.TotalReceived = transaction.Vout[i].Value
			address.TotalReceivedSat = address.TotalReceived * 100000000
			//address.TotalSent
			//address.TotalSentSat
			address.UnconfirmedBalance = 0
			address.UnconfirmedTxApperances = 0
			address.Transactions = append(address.Transactions, transaction.Txid)

			database.IndexAdress(db, address)
		}
		//transaction.Vout[i].ScriptPubKey.Addresses[0]
	}
}

// end result test 1 address
func test() {
	var addr address.Index
	decoder := gob.NewDecoder(bytes.NewReader(database.ViewAddress(db, "Ea6aiVS5dGWqVtQ4Akd9KCPw5HFmTbBPvX")))
	decoder.Decode(&addr)

	fmt.Printf("address: %s balance: %f", addr.AddrStr, addr.Balance)
	fmt.Println("\n")
	fmt.Println(addr)
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
