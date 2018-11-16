package rebuilddb

import (
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/go-errors/errors"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/insight"
	"github.com/romanornr/cyberchain/insightjson"
	"github.com/romanornr/cyberchain/mongodb"
	"github.com/romanornr/cyberchain/subsidy"
	"gopkg.in/cheggaaa/pb.v2"
	"io/ioutil"
	"log"
	"strings"
)

//var db = database.GetDatabaseInstance()

var db = mongodb.GetSession()
var isMainChain bool

func init() {
	ParseJson()
	IsMainChain()
}

func IsMainChain() {
	blockchainInfo, err := blockdata.GetBlockChainInfo()
	if err != nil {
		log.Fatalf("Error getting Blockchaininfo via RPC: %s", err)
	}

	if blockchainInfo.Chain != "main" {
		isMainChain = false
	}

	isMainChain = true
}

type Pools []struct {
	PoolName      string   `json:"poolName"`
	URL           string   `json:"url"`
	SearchStrings []string `json:"searchStrings"`
}

var pools Pools

// read and parse the json file and unmarshal
func ParseJson() {
	jsonFile, err := ioutil.ReadFile("../../pools.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(jsonFile), &pools)
}

/*  THIS WAS IN BBOLT/BOLTDB
note: 2000 blocks costs currently 8.4 MB and ~39 seconds to save. Running into performance issues.
goroutine "go addTransactions(block) speed up from ~39 seconds to ~29 seconds. 25% speed up

try to analyze this address: https://chainz.cryptoid.info/via/address.dws?369935.htm
*/

/*
	MongoDB
	2000 blocks without transactions cost 1.746 seconds
	2000 blocks with transactions cost 3.275 seconds

	2000 blocks with tx and a goroutine cost 2.56 seconds

    // for account balance check http://127.0.0.1:8000/api/via/addr/Vn5maEzzZNPQ85rKFAgACRW98oiDtmMumG
    // 6852caef331276d62c0de58ee430889c3926d9b5d832c7360dd9efe33fa1b6f6;11046;2014-07-20 00:04:49;10;2360.19357 <--block 11046
	// f5a38ecb879748de37c4bd4ae3695ae6fe324a61c666eae3d547e736ae42ff62;11129;2014-07-20 00:40:49;10;2540.20557
*/

func BuildDatabase() {
	//end := 	3673+200
	end := 11000 + 50
	//end := 11139 + 1
	progressBar := pb.StartNew(end)
	for i := 1; i < end; i++ {
		blockhash, _ := blockdata.GetBlockHash(int64(i))
		block, _ := blockdata.GetBlock(blockhash)
		newBlock, _ := insight.ConvertToInsightBlock(block)

		txs := GetTx(block)

		//add pool info to block before adding into mongodb
		coinbaseText := GetCoinbaseText(txs[0])
		pool, err := getPoolInfo(coinbaseText)
		if err == nil {
			newBlock.PoolInfo = &pool
		}

		//add reward info
		newBlock.Reward = subsidy.CalcViacoinBlockSubsidy(int32(newBlock.Height), isMainChain)
		newBlock.IsMainChain = isMainChain

		go mongodb.AddBlock(newBlock)

		AddTransactions(txs, newBlock.Height)

		progressBar.Increment()

	}
	progressBar.Finish()
}

// get coinbase hex string by getting the first transaction of the block
// in the tx.Vin[0] and decode the hex string into a normal text
// Example: "52062f503253482f04dee0c7530807ffffff010000000d2f6e6f64655374726174756d2f" -> /nodeStratum/
func GetCoinbaseText(tx *btcjson.TxRawResult) string {
	src := []byte(tx.Vin[0].Coinbase)

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}

	return string(dst[:n])
}

// range over all pools and within that range over all search strings
// check if a poolSearchString matches the coinbase text
func getPoolInfo(coinbaseText string) (insightjson.Pools, error) {
	var blockMinedByPool insightjson.Pools

	for _, pool := range pools {
		for _, PoolSearchString := range pool.SearchStrings {
			if strings.Contains(coinbaseText, PoolSearchString) {
				blockMinedByPool.PoolName = pool.PoolName
				blockMinedByPool.URL = pool.URL
				return blockMinedByPool, nil
			}
		}
	}
	return blockMinedByPool, errors.New("PoolSearchStrings did not match coinbase text. Unknown mining pool or solo miner")
}

func GetTx(block *btcjson.GetBlockVerboseResult) []*btcjson.TxRawResult {
	Transactions := []*btcjson.TxRawResult{}
	for i := 0; i < len(block.Tx); i++ {
		txhash, _ := chainhash.NewHashFromStr(block.Tx[i])
		tx, _ := blockdata.GetRawTransactionVerbose(txhash)
		Transactions = append(Transactions, tx)
	}

	return Transactions
}

func AddTransactions(transactions []*btcjson.TxRawResult, blockheight int64) {
	for _, transaction := range transactions {
		newTx := insight.TxConverter(transaction, blockheight)
		mongodb.AddTransaction(&newTx[0])
		CalcAddr(&newTx[0])
	}
}

func CalcAddr(tx *insightjson.Tx) {

	//receive
	for _, txVout := range tx.Vouts {

		for _, voutAdress := range txVout.ScriptPubKey.Addresses {
			dbAddrInfo, err := mongodb.GetAddressInfo(txVout.ScriptPubKey.Addresses[0])

			if err != nil {
				AddressInfo := insightjson.AddressInfo{
					voutAdress,
					txVout.Value,
					int64(txVout.Value * 100000000),
					txVout.Value,
					int64(txVout.Value * 100000000),
					0,
					0,
					0,
					0,
					0,
					1,
					[]string{tx.Txid},
				}
				mongodb.AddAddressInfo(&AddressInfo)
			} else {
				value := int64(txVout.Value * 100000000)
				mongodb.UpdateAddressInfoReceived(&dbAddrInfo, value, true, tx.Txid)
			}
		}

	}

	//sent
	for _, txVin := range tx.Vins {
		dbAddrInfo, err := mongodb.GetAddressInfo(txVin.Addr)
		value := int64(txVin.ValueSat)

		if err == nil {
			mongodb.UpdateAddressInfoSent(&dbAddrInfo, value, true, tx.Txid)
		}
	}
}
