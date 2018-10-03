package rebuilddb

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	_"github.com/romanornr/cyberchain/address"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"gopkg.in/cheggaaa/pb.v2"
	"log"
	"sync"
	"fmt"
	"github.com/romanornr/cyberchain/address"
	"encoding/gob"
	"bytes"
)

var db = database.GetDatabaseInstance()

/*
	note: 2000 blocks costs currently 8.4 MB and ~39 seconds to save. Running into performance issues.
	goroutine "go addTransactions(block) speed up from ~39 seconds to ~29 seconds. 25% speed up

	try to analyze this address: https://chainz.cryptoid.info/via/address.dws?369935.htm
*/


func BuildDatabase() {

	startblock := 5422072-1
	currentBlock := 5422072+10
	progressBar := pb.StartNew(currentBlock)
	wg := sync.WaitGroup{}
	wg.Add(currentBlock-startblock)

	for i := int64(startblock); i < int64(currentBlock); i++ {
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
			resolveAddresses(tx)
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

	        var addr address.Index

			addr.AddrStr = transaction.Vout[0].ScriptPubKey.Addresses[0]
			addr.TotalReceived = transaction.Vout[0].Value
			addr.TotalReceivedSat = addr.TotalReceived * 100000000
			//address.TotalSent
			//address.TotalSentSat
			addr.UnconfirmedBalance = 0
			addr.UnconfirmedTxApperances = 0
			addr.Transactions = append(addr.Transactions, transaction.Txid)

			//fmt.Println(addr.AddrStr)
			//fmt.Println(transaction.Txid)
			//fmt.Println(addr.TotalReceived)
			fmt.Println(transaction.Vin[0].ScriptSig)


			// Check if address was already in the database
			// use old values to calculate the new values.
			// addrInDatabase has all the values of what is already in the database
			var addrInDatabase address.Index
			addressInDatabase := database.ViewAddress(db, addr.AddrStr)
			if len(addressInDatabase) > 1 {
				decoder := gob.NewDecoder(bytes.NewReader(addressInDatabase))
				decoder.Decode(&addrInDatabase)

				addr.TotalReceived += addrInDatabase.TotalReceived
				addr.TotalReceivedSat += addrInDatabase.TotalReceivedSat
				addr.Transactions = append(addr.Transactions, addrInDatabase.Transactions[0])

				// delete old key in the database so the updated one can be inserted instead
				database.DeleteAddress(db, addr.AddrStr)

				}


			database.IndexAdress(db, addr)


	//var addr address.Index
	//
	//for i := 0; i < len(transaction.Vout); i++ {
	//	for j := 0; i < len(transaction.Vout[i].ScriptPubKey.Addresses); j++ {
	//		addr.AddrStr = transaction.Vout[i].ScriptPubKey.Addresses[j]
	//		addr.TotalReceived = transaction.Vout[i].Value
	//		addr.TotalReceivedSat = addr.TotalReceived * 100000000
	//		//address.TotalSent
	//		//address.TotalSentSat
	//		addr.UnconfirmedBalance = 0
	//		addr.UnconfirmedTxApperances = 0
	//		addr.Transactions = append(addr.Transactions, transaction.Txid)
	//
	//		fmt.Println(addr.AddrStr)
	//
	//		database.IndexAdress(db, addr)
	//	}
	//	//transaction.Vout[i].ScriptPubKey.Addresses[0]
	//}
}

//end result test 1 address
func test() {

	var addr address.Index
	fmt.Println(database.ViewAddress(db, "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A"))
	decoder := gob.NewDecoder(bytes.NewReader(database.ViewAddress(db, "VdMPvn7vUTSzbYjiMDs1jku9wAh1Ri2Y1A")))
	decoder.Decode(&addr)

	fmt.Printf("address: %s received: %f\n", addr.AddrStr, addr.TotalReceived)
	//fmt.Println(addr)
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
