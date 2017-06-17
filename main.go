package main

import (
	"log"

	"github.com/btcsuite/btcrpcclient"
)

func main() {
	// Connect to local RPC server using HTTP POST mode.
	connCfg := &btcrpcclient.ConnConfig{
		Host:         "localhost:8332",
		User:         "viarpc",
		Pass:         "viapass",
		HTTPPostMode: true, // Bitcoin Core & Bitcoin based altcoins only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core $ Bitcoin based altcoins do not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := btcrpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
}