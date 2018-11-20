package zeroMQ

import (
	"fmt"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/mongodb"
	"github.com/romanornr/cyberchain/notification"
	"github.com/zeromq/goczmq"
	"log"
	"github.com/spf13/viper"
)

func BlockNotify() {

	endpoint := viper.GetString("zmq.endpoint")

	subscriber, err := goczmq.NewSub(endpoint, "hashblock")
	defer subscriber.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ZeroMQ started to listen for blocks")

	for {
		msg, _, err := subscriber.RecvFrame()

		if err != nil {
			log.Printf("Error ZMQ RecFrame: %s", err)
		}

		//lenght of a hash
		if len(msg) == 32 {
			block, err := blockdata.GetLatestBlock()
			if err != nil {
				log.Printf("Error ZMQ getting latest block: %s", err)
			}

			lastBlockInDb, err := mongodb.GetLastBlock()
			if lastBlockInDb.Height+int64(1) != block.Height {
				log.Printf("Warning: new block via ZMQ has blockheight %d but last blockheight in DB is %d\n", block.Height, lastBlockInDb.Height)
				log.Print("It could be that all blocks in MongoDB are not up-to-date.")
			}

			notification.ProcessBlock(block)
			fmt.Println("block added")
		}

	}

}
