// Copyright (c) 2019 Romano, Viacoin developer
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rebuilddb

import (
	"github.com/romanornr/blockexplorer/blockdata"
	"github.com/romanornr/blockexplorer/notification"
	"gopkg.in/cheggaaa/pb.v2"
)

/*  THIS WAS IN BBOLT/BOLTDB
note: 2000 blocks costs currently 8.4 MB and ~39 seconds to save. Running into performance issues.
goroutine "go addTransactions(block) speed up from ~39 seconds to ~29 seconds. 25% speed up

try to analyze this address: https://chainz.cryptoid.info/via/address.dws?369935.htm
*/

/*
	MongoDB
	2000 blocks without transactions cost 1.746 seconds
	2000 blocks with transactions cost 3.275 seconds

	2000 blocks with tx and a goroutine cost 2.56 seconds

    // for account balance check http://127.0.0.1:8000/api/via/addr/Vn5maEzzZNPQ85rKFAgACRW98oiDtmMumG
    // 6852caef331276d62c0de58ee430889c3926d9b5d832c7360dd9efe33fa1b6f6;11046;2014-07-20 00:04:49;10;2360.19357 <--block 11046
	// f5a38ecb879748de37c4bd4ae3695ae6fe324a61c666eae3d547e736ae42ff62;11129;2014-07-20 00:40:49;10;2540.20557
*/

// receive the latest blockheight
// RPC call and add every block into the database
// from 1 till the latestBlockHeight
func BuildDatabase(latestBlockHeight int64) {
	progressBar := pb.Start64(latestBlockHeight)
	for i := int64(1); i < latestBlockHeight; i++ {
		blockHash, _ := blockdata.GetBlockHash(i)
		block, _ := blockdata.GetBlock(blockHash)
		notification.ProcessBlock(block)
		progressBar.Increment()

	}
	progressBar.Finish()
}
