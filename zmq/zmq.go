package zeroMQ

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/blockexplorer/blockdata"
	"github.com/romanornr/blockexplorer/mongodb"
	"github.com/romanornr/blockexplorer/notification"
	"github.com/spf13/viper"
	"github.com/zeromq/goczmq"
)

var dao = mongodb.MongoDAO{
	"127.0.0.1",
	"viacoin",
}

// listen to ZMQ endpoint and check for message length if it's a hash
// When a new Block is available the function isSynced() will be called
// isSynced() return the new block from RPC and an error message
// the error message means the block explorer is not in sync and should add the block
// received from isSynced() to be insync.
func BlockNotify() {

	endpoint := viper.GetString("zmq.endpoint")

	subscriber, err := goczmq.NewSub(endpoint, "hashblock")
	defer subscriber.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("ZeroMQ started to listen for blocks")

	for {
		msg, _, err := subscriber.RecvFrame()

		if err != nil {
			log.Warningf("Error ZMQ RecFrame: %s", err)
		}

		//lenght of a hash
		if len(msg) == 32 {
			block, _, blocksBehind, _ := isSyned()
			//if err != nil {
				go notification.ProcessBlock(block)
			//}
			log.WithFields(log.Fields{
				"block height": block.Height,
				"block hash":   block.Hash,
			}).Info("block successfully added to the database")
			log.Warningf("blocks behind: %d\n", blocksBehind)
		}

	}
}

// check if block explorer is synced
// meaning the last block in the database should be equal to the block received from an RPC call
// always return the best block (block received from RPC call) how many blocks it's behind and error message
// being behind 1 block behind is not an issue. Just needs the newest block
func isSyned() (bestBlock *btcjson.GetBlockVerboseResult, synced bool, blocksBehind int64, err error) {
	block, err := blockdata.GetLatestBlock()
	if err != nil {
		log.Warningf("RPC call failed to get latest block: %s\n", err)
		return block, false, int64(0), fmt.Errorf("RPC call failed to get latest block: %s\n", err)
	}

	lastBlockInDb, err := dao.GetLastBlock()
	if err != nil {
		return block, false, int64(0), fmt.Errorf("failed to retrieve last block from database: %s\n", err)
	}

	blocksBehind = block.Height - lastBlockInDb.Height

	if lastBlockInDb.Height != block.Height {
		return block, false, blocksBehind, nil
	}

	if lastBlockInDb.Hash != block.Hash {
		return block, false, blocksBehind, nil
	}

	return block, true, blocksBehind, nil
}
