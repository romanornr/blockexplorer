package blockdata

import (
	"log"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/spf13/viper"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func client() *rpcclient.Client {
	// Connect to local bitcoin/altcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         viper.GetString("rpc.ip") + ":" + viper.GetString("rpc.port"), //127.0.0.1:8332
		User:         viper.GetString("rpc.username"),
		Pass:         viper.GetString("rpc.password"),
		HTTPPostMode: true, // Viacoin core only supports HTTP POST mode
		DisableTLS:   true, // Viacoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
		client.Shutdown()
	}
	//defer client.Shutdown()

	return client
}

// get current difficulty of a block
func GetDifficulty() float64{
	difficulty, err := client().GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	return difficulty
}

func GetBlockCount() int64{
	c, err := client().GetBlockCount()
	if err != nil{
		log.Fatal(err)
	}
	return c
}

func GetBlockHash(blockHeight int64) *chainhash.Hash  {
	h, err := client().GetBlockHash(blockHeight)
	if err != nil{
		log.Fatal(err)
	}
	return h
}

func GetBlock(blockHash *chainhash.Hash) *wire.MsgBlock{
	block, err := client().GetBlock(blockHash)
	if err != nil{
		log.Fatal(err)
	}
	return block
}

// get latest block info
func GetLatestBlockInfo() *btcjson.GetBlockVerboseResult{
	blockCount, err := client().GetBlockCount() //get the latest blocks
	if err != nil {
		log.Println(err)
	}
	hash, err := client().GetBlockHash(blockCount)
	if err != nil {
		log.Println(err)
	}

	block, err := client().GetBlockVerbose(hash)
	if err != nil {
		log.Fatal(err)
	}

	return block
}