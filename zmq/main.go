package main

import (
	"fmt"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/mongodb"
	"github.com/romanornr/cyberchain/notification"
	"github.com/zeromq/goczmq"
	"log"
)

func ZeroMQBlockNotify() {

	mongodb.DropDatabase()

	subscriber, err := goczmq.NewSub("tcp://127.0.0.1:28332", "hashblock")
	if err != nil {
		log.Fatal(err)
	}

	defer subscriber.Destroy()

	for {
		msg, _, err := subscriber.RecvFrame()

		if err != nil {
			log.Printf("Error ZMQ RecFrame: %s", err)
		}

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
