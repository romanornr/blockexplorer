// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package Reorg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/mongodb"
)

func Check(newBlock *btcjson.GetBlockVerboseResult) {
	tip, _ := mongodb.GetLastBlock()
	if newBlock.PreviousHash != tip.Hash {
		log.Warningln("Reorg detected")
	}
}
