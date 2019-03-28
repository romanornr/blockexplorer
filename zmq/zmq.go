package zeroMQ

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/mongodb"
	"github.com/romanornr/cyberchain/notification"
	"github.com/spf13/viper"
	"github.com/zeromq/goczmq"
)

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
			block, err := isSyned()
			if err != nil {
				go notification.ProcessBlock(block)
			}
			log.WithFields(log.Fields{
				"block height": block.Height,
				"block hash":   block.Hash,
			}).Info("block successfully added to the database")
		}

	}
}

// check if block explorer is synced
// meaning the last block in the database should be equal to the block received from an RPC call
// always return the best block (block received from RPC call) and error message
func isSyned() (bestBlock *btcjson.GetBlockVerboseResult, err error) {
	block, err := blockdata.GetLatestBlock()
	if err != nil {
		log.Warningf("RPC call failed to get latest block: %s\n", err)
		return block, fmt.Errorf("RPC call failed to get latest block: %s\n", err)
	}

	lastBlockInDb, err := mongodb.GetLastBlock()
	if lastBlockInDb.Hash != block.Hash {
		return block, fmt.Errorf("block explorer is not in sync. Blockhash in database: %s Blockhash via RPC call: %s", lastBlockInDb.Hash, block.Hash)
	}

	if lastBlockInDb.Height != block.Height {
		return block, fmt.Errorf("block explorer is not in sync. BlockHeight in database: %d Blockheight via RPC call: %d", lastBlockInDb.Height, block.Height)
	}
	return block, nil
}
