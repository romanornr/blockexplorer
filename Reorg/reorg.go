// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package Reorg

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/blockexplorer/insightjson"
	"github.com/romanornr/blockexplorer/mongodb"
)

func Check(dao mongodb.MongoDAO, block *btcjson.GetBlockVerboseResult) (reorg bool, tip insightjson.BlockResult, newBlock *btcjson.GetBlockVerboseResult) {
	tip, err := dao.GetLastBlock()
	if err != nil {
		log.Warningf("error getting tip from database: %s\n", err)
	}
	if block.PreviousHash != tip.Hash {
		log.Warningln("Reorg detected")
		return true, tip, block
	}

	return false, tip, block
}

func RollbackTransaction() {

}

func RollbackAddrIndex(dao mongodb.MongoDAO, tx *insightjson.Tx) {
	fmt.Println(tx.Vouts)
	//receive
	for _, txVout := range tx.Vouts {
		for _, voutAdress := range txVout.ScriptPubKey.Addresses {
			dbAddrInfo, _ := dao.GetAddressInfo(txVout.ScriptPubKey.Addresses[0])
			value := int64(txVout.Value * 100000000) // satoshi value to coin value
			log.Infof("rolling back address info for %s\n", voutAdress)
			dao.RollbackAddressInfoReceived(&dbAddrInfo, value, true, tx.Txid)
		}
	}

	//sent
	//for _, txVin := range tx.Vins {  //TODO ROLLBACK SENT
	//
	//	dbAddrInfo, err := dao.GetAddressInfo(txVin.Addr)
	//	value := int64(txVin.ValueSat)
	//
	//	if err == nil {
	//		go dao.UpdateAddressInfoSent(&dbAddrInfo, value, true, tx.Txid)
	//	}
	//}
}
