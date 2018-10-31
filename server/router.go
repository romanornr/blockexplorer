package server

import (
	_"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/julienschmidt/httprouter"
	"github.com/romanornr/cyberchain/Blockchain"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"github.com/romanornr/cyberchain/insight"
	"github.com/spf13/viper"
	"fmt"
)

var db = database.GetDatabaseInstance()

// createRouter creates and returns a router.
func createRouter() *httprouter.Router {
	network := viper.GetString("coin.symbol")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/api/"+network+"/getdifficulty", getDifficulty)
	router.GET("/api/"+network+"/blocks", getLatestBlocks)
	router.GET("/api/"+network+"/block/:hash", getBlock)
	router.GET("/api/"+network+"/block-index/:height", getBlockIndex)
	router.GET("/api/"+network+"/tx/:txid", getTransaction)

	fileServer := http.FileServer(http.Dir("static"))

	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		fileServer.ServeHTTP(w, r)
	})

	return router
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	coin := viper.Get("coin.name")
	err := tpl.ExecuteTemplate(w, "index.gohtml", coin)
	if err != nil {
		log.Println("error")
	}
}

func getDifficulty(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	difficulty, err := blockdata.GetDifficulty()
	if err != nil {
		log.Println("Could not get difficulty", err)
	}
	json.NewEncoder(w).Encode(difficulty)
}

// getLatestBlocks gets x (int) latest blocks
func getLatestBlocks(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	blockCount := blockdata.GetBlockCount() // get the latest blocks

	var blocks []*btcjson.GetBlockVerboseResult

	// blockheight - 1 in the loop. Get the blockhash from the height
	for i := 0; i < 10; i++ {
		prevBlock := blockCount - int64(i)
		hash := blockdata.GetBlockHash(prevBlock)

		block, err := blockdata.GetBlock(hash)
		if err != nil {
			log.Fatal(err)
		}

		blocks = append(blocks, block)
	}
	json.NewEncoder(w).Encode(blocks)
}

// get block by blockhash using the api url
// example: http://127.0.0.1:8000/api/via/block/45c2eb3f3ca602e36b9fac0c540cf2756f1d41719b4be25adb013f87bafee7bc
func getBlock(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	hash, err := chainhash.NewHashFromStr(ps.ByName("hash"))
	if err != nil {
		log.Printf("could not convert string to hash: %s\n", err)
	}

	proxy := Blockchain.BlockListProxy{}

	block, err := proxy.FindBlock(hash)
	if err != nil {
		log.Printf("error finding block: %s", err)
	}

	block.Confirmations = getBlockConfirmations(*block) // needs dynamic calculation

	apiblock, err := insight.ConvertToInsightBlock(block)

	json.NewEncoder(w).Encode(&apiblock)
}

// confirmations from blocks always change. Block confirmations can be calculated with the following method
// latest blockheight in database - blockheight
func getBlockConfirmations(block btcjson.GetBlockVerboseResult) int64 {
	currentBlockHeight, _ := database.GetLastBlockHeight(db)

	blockConfirmations := int64(currentBlockHeight) - block.Confirmations
	var m int64 = blockConfirmations
	var q = &m
	block.Confirmations = *q
	return *q
}

func getBlockIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	height := ps.ByName("height")
	blockheight, err := strconv.ParseUint(height, 10, 64)
	if err != nil {
		log.Println("could not convert height to int64")
	}

	proxy := Blockchain.BlockListProxy{}
	blockIndex := proxy.FindBlockHash(int64(blockheight))

	json.NewEncoder(w).Encode(blockIndex)
}

func getTransaction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	txid := ps.ByName("txid")
	txhash,_ := chainhash.NewHashFromStr(txid)
	//var tx *insightjson.Tx
	//decoder := gob.NewDecoder(bytes.NewReader(database.GetTransaction(db, txid)))
	//decoder.Decode(&tx)
	//
	//json.NewEncoder(w).Encode(tx)
	tx := blockdata.GetRawTransactionVerbose(txhash)
	txnew := insight.TxConverter(tx)
	fmt.Println(txnew[0].ValueOut)
	json.NewEncoder(w).Encode(txnew[0])

}
