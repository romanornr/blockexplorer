// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package Reorg

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/romanornr/blockexplorer/insightjson"
	"github.com/romanornr/blockexplorer/mongodb"
	"github.com/romanornr/blockexplorer/notification"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/blockexplorer/blockdata"
)

var dao = mongodb.MongoDAO{
	"localhost",
	"viacoin",
}

func BuildMockDatabase() {
	dao.DropDatabase()
	for i := int64(1); i < 6; i++ {
		blockhash, _ := blockdata.GetBlockHash(i)
		block, _ := blockdata.GetBlock(blockhash)
		notification.ProcessBlock(block)
	}
}

type fakeBlocks struct {
	height int64
	hash   string
}

var hashes = [7]string{"5ca83af67146e286610e118cc8f8e6a183c319fbb4a8fdb9e99daa2b8a29b3e3",
	"45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc",
	"7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2",
	"a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f",
	"d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6", //5
	"e8957dac3477849c431dce6929e45ca829598bf45f05f776742f04f06c246ae7",
	"5ca78b039ccfec56373a4392c043bb9a6c77f8c2934af96b036c00dd2e4a0cfa",
}

// this block has a blockheight:2 but this blockheight is already in the database
// It has a different hash compared to the block in the database with blockheight:2
// This would be a potential chain reorg
var reorgBlock = &btcjson.GetBlockVerboseResult{
	Hash:         "8b1419de52400f6467d311c9d6a5e4fd8a0816041fb7572ffc704fd7f9ffe8ef",
	PreviousHash: "a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f",
	Height:       4,
}

// block d8c9053f3c807b1465bd0a8bc99421e294066dd59e98cf14bb49d990ea88aff6
// has previous hash: a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f
// we insert a new block with previous hash a35d1bdbd41ea6c290d9a151bdafd39b76eda3c9c9d44e02d0209dd77f5aec1f
// in order to trigger a reorg.
func TestComparePreviousHash(t *testing.T) {

	dao.Connect()
	BuildMockDatabase()

	blockhash, err := chainhash.NewHashFromStr(hashes[5]) // check if blockhash is valid
	if err != nil {
		t.Errorf("The hash %s is not valid\n", blockhash.String())
	}

	block, err := blockdata.GetBlock(blockhash)
	if err != nil {
		t.Errorf("Could not get block: %s via RPC", hashes[5])
	}

	reorg, _, _ := Check(dao, block)

	if reorg != false {
		t.Errorf("Did not expect reorg: %s", err)
	}

	//check if it detects a chain reorg
	//blockhash, _ = chainhash.NewHashFromStr(hashes[2]) // check if blockhash is valid
	//block, _ = blockdata.GetBlock(blockhash)
	reorg, _, _ = Check(dao, reorgBlock)
	if reorg != true {
		t.Errorf("No chain reorg detected, however it was exected")
	}
}

var tables = []struct {
	blockHeight int64
}{
	{blockHeight: 5422072},
	{blockHeight: 5422822},
	{blockHeight: 5423010},
	{blockHeight: 5425232},
	{blockHeight: 5455537},
	{blockHeight: 5473972},
	{blockHeight: 5590354},
	{blockHeight: 5671670},
	{blockHeight: 5891328},
	{blockHeight: 6075611},
	{blockHeight: 6075728},
	{blockHeight: 6075872},
	{blockHeight: 6142214},
}

func TestRollbackTipTransaction(t *testing.T) {

	dao.DropDatabase()

	//adding blocks
	for _, table := range tables {
		blockHash, _ := blockdata.GetBlockHash(table.blockHeight)
		block, _ := blockdata.GetBlock(blockHash)
		notification.ProcessBlock(block)
	}

	addr, err := dao.GetAddressInfo("Ea6aiVS5dGWqVtQ4Akd9KCPw5HFmTbBPvX")
	if err != nil {
		fmt.Println(err)
	}

	logrus.Infof("address balance before rollbacks: %.8f\n", addr.Balance)

	if addr.BalanceSat != 534180377091241 {
		t.Errorf("Expected: %f, actual: %f\n",534180377091241, addr.Balance)
	}

	// rolling back block 6142214
	blockHash, _ := blockdata.GetBlockHash(6142214)
	block, _ := blockdata.GetBlock(blockHash)

	btcjsonTransactions := notification.GetTx(block)

	var transactions []insightjson.Tx
	// get insightjson format transaction from database
	for _, rawTransaction := range btcjsonTransactions {
		tx, _ := dao.GetTransaction(rawTransaction.Txid)
		transactions = append(transactions, tx)
	}

	for _, transaction := range transactions {
		RollbackAddrIndex(dao, &transaction)
	}

	addr, err = dao.GetAddressInfo("Ea6aiVS5dGWqVtQ4Akd9KCPw5HFmTbBPvX")
	if err != nil {
		fmt.Println(err)
	}

	logrus.Infof("address balance after rolling back tip: %f\n", addr.Balance)
}

//// Here we use this address for rollbacks
//// https://chainz.cryptoid.info/via/address.dws?369935.htm
//func TestRollbackTransaction(t *testing.T) {
//
//	dao.DropDatabase()
//
//	//adding blocks
//	for _, table := range tables {
//		blockHash, _ := blockdata.GetBlockHash(table.blockHeight)
//		block, _ := blockdata.GetBlock(blockHash)
//		notification.ProcessBlock(block)
//	}
//
//	var transactions []insightjson.Tx
//
//	// getting transactions
//	for _, table := range tables {
//		blockHash, _ := blockdata.GetBlockHash(table.blockHeight)
//		block, _ := blockdata.GetBlock(blockHash)
//		btcjsonTransactions := notification.GetTx(block)
//
//		// get insightjson format transaction from database
//		for _, rawTransaction := range btcjsonTransactions {
//			tx, _ := dao.GetTransaction(rawTransaction.Txid)
//			transactions = append(transactions, tx)
//		}
//	}
//
//	for _, transaction := range transactions {
//		RollbackAddrIndex(dao, &transaction)
//	}
//
//	//txhash, _ := chainhash.NewHashFromStr("")
//	//txdb, _ := dao.GetTransaction(*txhash)
//	//fmt.Println(txdb.Vouts)
//	//
//	//RollbackAddrIndex(dao, &txdb)
//
//	addr, err := dao.GetAddressInfo("Ea6aiVS5dGWqVtQ4Akd9KCPw5HFmTbBPvX")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Printf("adress balance is: %f\n", addr.Balance)
//}
