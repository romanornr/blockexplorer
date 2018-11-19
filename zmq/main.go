package main

import (
	"github.com/zeromq/goczmq"
	"log"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/notification"
	"fmt"
	"github.com/romanornr/cyberchain/mongodb"
)


func main() {

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

			notification.ProcessBlock(block)
			fmt.Println("block added")
		}

	}

}

func ProcessBlock() {


}