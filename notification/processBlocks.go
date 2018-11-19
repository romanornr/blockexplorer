package notification

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/insight"
	"github.com/romanornr/cyberchain/mongodb"
	"github.com/romanornr/cyberchain/insightjson"
	"encoding/hex"
	"log"
	"strings"
	"github.com/romanornr/cyberchain/blockdata"
	"io/ioutil"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/romanornr/cyberchain/subsidy"
	"runtime"
	"path/filepath"
)

var db = mongodb.GetSession()
var isMainChain bool

func init() {
	parseJson()
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

// get the current path: notification/
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// read and parse the json file and unmarshal
func parseJson() {
	path := strings.Split(basepath, "notification")
	jsonFile, err := ioutil.ReadFile(path[0] + "pools.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(jsonFile), &pools)
}

func ProcessBlock(block *btcjson.GetBlockVerboseResult) {
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

	AddTransactions(txs, newBlock.Height) // this in a go routine def causes a race conditions
}


// get coinbase hex string by getting the first transaction of the block
// in the tx.Vin[0] and decode the hex string into a normal text
// Example: "52062f503253482f04dee0c7530807ffffff010000000d2f6e6f64655374726174756d2f" -> /nodeStratum/
func GetCoinbaseText(tx *btcjson.TxRawResult) string {
	src := []byte(tx.Vin[0].Coinbase)

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Printf("Error getting coinbase text: %s", err)
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
		go mongodb.AddTransaction(&newTx[0])
		AddrIndex(&newTx[0]) //this in a go routine will cause a race condition
	}
}

func AddrIndex(tx *insightjson.Tx) {
	//receive
	for _, txVout := range tx.Vouts {
		go func() {
			for _, voutAdress := range txVout.ScriptPubKey.Addresses {
				dbAddrInfo, err := mongodb.GetAddressInfo(txVout.ScriptPubKey.Addresses[0])
				if err != nil {
					addressInfo := createAddressInfo(voutAdress, txVout, tx)
					go mongodb.AddAddressInfo(addressInfo)
				} else {
					value := int64(txVout.Value * 100000000) // satoshi value to coin value
					go mongodb.UpdateAddressInfoReceived(&dbAddrInfo, value, true, tx.Txid)
				}
			}
		}()
	}

	//sent
	for _, txVin := range tx.Vins {

		dbAddrInfo, err := mongodb.GetAddressInfo(txVin.Addr)
		value := int64(txVin.ValueSat)

		if err == nil {
			go mongodb.UpdateAddressInfoSent(&dbAddrInfo, value, true, tx.Txid)
		}
	}
}

// create address info. An address can only "exist" if it ever received a transaction
// the received is the vout values.
func createAddressInfo(address string, txVout *insightjson.Vout, tx *insightjson.Tx) *insightjson.AddressInfo {
	addressInfo := insightjson.AddressInfo{
		address,
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
	return &addressInfo
}