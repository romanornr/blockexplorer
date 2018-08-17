package blockdata

import (
	"log"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	client2 "github.com/romanornr/cyberchain/client"
)

func client() *rpcclient.Client {
	// Connect to local bitcoin/altcoin core RPC server using HTTP POST mode.
	//connCfg := &rpcclient.ConnConfig{
	//	Host:         "127.0.0.1:5222",
	//	User:         "via",
	//	Pass:         "via",
	//	HTTPPostMode: true, // Viacoin core only supports HTTP POST mode
	//	DisableTLS:   true, // Viacoin core does not provide TLS by default
	//}

	connCfg := client2.LoadConfig()c,
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


func GetBlock(blockhash *chainhash.Hash) *btcjson.GetBlockVerboseResult {
	block, err := client().GetBlockVerboseTx(blockhash)
	if err != nil {
		log.Fatal(err)
	}
	return block
}

//func GetBlock(blockhash *chainhash.Hash) *wire.MsgBlock {
//	block, err := client().GetBlock(blockhash)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return block
//}

func GetBlockHeader(blockhash *chainhash.Hash) *btcjson.GetBlockHeaderVerboseResult{
	block, err := client().GetBlockHeaderVerbose(blockhash)
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

func GetRawTransactionVerbose(transactionHash *chainhash.Hash) *btcjson.TxRawResult {
	rawtx, err := client().GetRawTransactionVerbose(transactionHash)
	if err != nil {
		log.Println(err)
	}

	return rawtx
}

func GetRawTransaction(transactionHash *chainhash.Hash) *btcutil.Tx{
	rawtx, err := client().GetRawTransaction(transactionHash)
	if err != nil {
		log.Println(err)
	}
	return rawtx
}

// Decode the raw transaction hash into a human readable json
func DecodeRawTransaction(transactionHash []byte) *btcjson.TxRawResult{
	decodedRawTransaction, err := client().DecodeRawTransaction(transactionHash)
	if err != nil {
		log.Println(err)
	}
	return decodedRawTransaction
}