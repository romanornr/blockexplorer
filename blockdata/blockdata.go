package blockdata

import (
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/romanornr/cyberchain/client"
)

var rpclient = client.GetInstance()

func GetBlockCount() int64 {
	c, err := rpclient.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func GetBlockHash(blockHeight int64) (*chainhash.Hash, error) {
	blockhash, err := rpclient.GetBlockHash(blockHeight)
	return blockhash, err
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

func GetBlockHeader(blockhash *chainhash.Hash) (*btcjson.GetBlockHeaderVerboseResult, error) {
	block, err := rpclient.GetBlockHeaderVerbose(blockhash)
	return block, err
}

func GetBlockHeaderVerbose(blockhash *chainhash.Hash) (*btcjson.GetBlockHeaderVerboseResult, error) {
	block, err := rpclient.GetBlockHeaderVerbose(blockhash)
	return block, err
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

func GetRawTransactionVerbose(transactionHash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	rawtx, err := rpclient.GetRawTransactionVerbose(transactionHash)
	return rawtx, err
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

// get current difficulty of a block
func GetDifficulty() (float64, error) {
	difficulty, err := rpclient.GetDifficulty()
	return difficulty, err
}

func SearchRawTransactionsVerbose(address btcutil.Address, skip, count int, reverse bool, filterAddrs []string) []*btcjson.SearchRawTransactionsResult {
	tx, err := rpclient.SearchRawTransactionsVerbose(address, skip, count, reverse, reverse, filterAddrs)
	if err != nil {
		log.Printf("Search raw tx error: %s", err)
	}
	return tx
}


func GetBlockChainInfo() (*btcjson.GetBlockChainInfoResult, error) {
	result, err := rpclient.GetBlockChainInfo()
	return result, err
}