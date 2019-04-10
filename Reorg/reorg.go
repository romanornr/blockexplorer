// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package Reorg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/blockexplorer/insightjson"
	"github.com/romanornr/blockexplorer/mongodb"
)

func Check(dao mongodb.MongoDAO, block *btcjson.GetBlockVerboseResult) (reorg bool, tip insightjson.BlockResult,  newBlock *btcjson.GetBlockVerboseResult) {
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
