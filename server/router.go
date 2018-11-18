package server

import (
	_ "bytes"
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/julienschmidt/httprouter"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
	"github.com/romanornr/cyberchain/mongodb"
	"fmt"
)

var db = database.GetDatabaseInstance()
var coin = viper.Get("coin.name")

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
	router.GET("/api/"+network+"/addr/:addr", getAddressInfo)

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

	err := tpl.ExecuteTemplate(w, "index.gohtml", coin)
	if err != nil {
		log.Printf("Error executing template: %s", err)
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
		hash, _ := blockdata.GetBlockHash(prevBlock)

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

	//proxy := Blockchain.BlockListProxy{}
	//
	//block, err := proxy.FindBlock(hash)
	//if err != nil {
	//	log.Printf("error finding block: %s", err)
	//}

	block, err := mongodb.GetBlock(*hash)
	if err != nil {
		fmt.Fprintf(w, "Not found")
		return
	}
	//block.Confirmations = getBlockConfirmations(block)
	//block.Confirmations = getBlockConfirmations(*block) // needs dynamic calculation

	//apiblock, err := insight.ConvertToInsightBlock(block)

	json.NewEncoder(w).Encode(&block)
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

	//proxy := Blockchain.BlockListProxy{}
	//blockIndex := proxy.FindBlockHash(int64(blockheight))
	blockIndex, err := blockdata.GetBlockHash(int64(blockheight))
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	json.NewEncoder(w).Encode(blockIndex)
}

func getTransaction(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	txid := ps.ByName("txid")
	txhash, _ := chainhash.NewHashFromStr(txid)

	tx, err := mongodb.GetTransaction(*txhash)
	if err != nil {
		fmt.Fprintf(w, "Not found")
		return
	}
	json.NewEncoder(w).Encode(tx)
}

func getAddressInfo(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	address := ps.ByName("addr")

	addressInfo, err := mongodb.GetAddressInfo(address)
	if err != nil {
		fmt.Println("shit")
	}

	json.NewEncoder(w).Encode(addressInfo)
}