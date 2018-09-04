package server

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/romanornr/cyberchain/blockdata"
	"github.com/romanornr/cyberchain/database"
	"strconv"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"encoding/gob"
	"bytes"
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

func getBlock(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	hash, err := chainhash.NewHashFromStr(ps.ByName("hash"))
	if err != nil {
		log.Println(err)
	}

	blockInDatabase := database.ViewBlock(db, hash.String())

	var block *btcjson.GetBlockVerboseResult

	decoder := gob.NewDecoder(bytes.NewReader(blockInDatabase))
	decoder.Decode(&block)

	block.Confirmations = getBlockConfirmations(*block) // needs dynamic calculation

	json.NewEncoder(w).Encode(&block)
}

// confirmations from blocks always change. Block confirmations can be calculated with the following method
// latest blockheight in database - blockheight
func getBlockConfirmations(block btcjson.GetBlockVerboseResult) int64{
	currentBlockHeight, _ := database.GetLastBlockHeight(db)

	blockConfirmations := int64(currentBlockHeight) - block.Confirmations
	var m int64 = blockConfirmations
	var q = &m
	block.Confirmations = *q
	return *q
}

type blockIndex struct {
	BlockHash string `json:"blockHash"`
}

func getBlockIndex(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	height := ps.ByName("height")
	i64, err := strconv.ParseUint(height, 10, 64)
	if err != nil {
		log.Println("could not convert height to int64")
	}

	hash := string(database.FetchBlockHashByBlockHeight(int64(i64)))

	json.NewEncoder(w).Encode(blockIndex{hash})
}