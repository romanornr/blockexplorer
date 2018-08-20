package blockdata

import (
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/romanornr/cyberchain/client"
)

var rpclient = client.GetInstance()

// get current difficulty of a block
func GetDifficulty() float64 {
	difficulty, err := rpclient.GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	return difficulty
}

func GetBlockCount() int64 {
	c, err := rpclient.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func GetBlockHash(blockHeight int64) *chainhash.Hash {
	h, err := rpclient.GetBlockHash(blockHeight)
	if err != nil {
		log.Fatal(err)
	}
	return h
}

func GetBlockHashAsync(blockHeight int64) *chainhash.Hash {
	f, err := rpclient.GetBlockHashAsync(blockHeight).Receive()
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func GetBlock(blockhash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	// block, err := rpclient.GetBlockVerboseTx(blockhash)
	block, err := rpclient.GetBlockVerbose(blockhash)
	// if err != nil {
	// 	log.Fatalf("Block with hash: %s: %s\n", blockhash.String(), err)
	// }
	return block, err
}

func GetBlockAsync(blockhash *chainhash.Hash) *btcjson.GetBlockVerboseResult {
	block, err := rpclient.GetBlockVerboseAsync(blockhash).Receive()
	if err != nil {
		log.Fatal(err)
	}
	return block
}

// func GetBlock(blockhash *chainhash.Hash) *wire.MsgBlock {
// 	block, err := client().GetBlock(blockhash)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return block
// }

func GetBlockHeader(blockhash *chainhash.Hash) *btcjson.GetBlockHeaderVerboseResult {
	block, err := rpclient.GetBlockHeaderVerbose(blockhash)
	if err != nil {
		log.Fatal(err)
	}
	return block
}

// get latest block info
func GetLatestBlockInfo() *btcjson.GetBlockVerboseResult {
	blockCount, err := rpclient.GetBlockCount() // get the latest blocks
	if err != nil {
		log.Println(err)
	}
	hash, err := rpclient.GetBlockHash(blockCount)
	if err != nil {
		log.Println(err)
	}

	block, err := rpclient.GetBlockVerbose(hash)
	if err != nil {
		log.Fatal(err)
	}

	return block
}

func GetRawTransactionVerbose(transactionHash *chainhash.Hash) *btcjson.TxRawResult {
	rawtx, err := rpclient.GetRawTransactionVerbose(transactionHash)
	if err != nil {
		log.Println(err)
	}

	return rawtx
}

func GetRawTransaction(transactionHash *chainhash.Hash) *btcutil.Tx {
	rawtx, err := rpclient.GetRawTransaction(transactionHash)
	if err != nil {
		log.Println(err)
	}
	return rawtx
}

// Decode the raw transaction hash into a human readable json
func DecodeRawTransaction(transactionHash []byte) *btcjson.TxRawResult {
	decodedRawTransaction, err := rpclient.DecodeRawTransaction(transactionHash)
	if err != nil {
		log.Println(err)
	}
	return decodedRawTransaction
}
