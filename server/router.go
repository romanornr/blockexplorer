package server

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/romanornr/cyberchain/database"
	"github.com/spf13/viper"
	"github.com/viacoin/viad/btcjson"
)

// createRouter creates and returns a router.
func createRouter() *httprouter.Router {
	network := viper.GetString("coin.symbol")

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/api/"+network+"/getdifficulty", getDifficulty)
	router.GET("/api/"+network+"/blocks", getLatestBlocks)
	router.GET("/api/"+network+"/block", getBlock)

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
	difficulty, err := client().GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(difficulty)
}

// getLatestBlocks gets x (int) latest blocks
func getLatestBlocks(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	blockCount, err := client().GetBlockCount() // get the latest blocks
	if err != nil {
		log.Println(err)
	}

	var blocks []*btcjson.GetBlockVerboseResult

	// blockheight - 1 in the loop. Get the blockhash from the height
	// 10? why 10?
	for i := 0; i < 10; i++ {
		prevBlock := blockCount - int64(i)
		hash, err := client().GetBlockHash(prevBlock)
		if err != nil {
			log.Println(err)
		}

		block, err := client().GetBlockVerbose(hash)
		if err != nil {
			log.Fatal(err)
		}

		blocks = append(blocks, block)

	}
	json.NewEncoder(w).Encode(blocks)
}

func getBlock(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	x := database.ViewBlock("c65aabe1578da37945118ff6078792315499dd0dd6712f76a5f387126799d9b1")

	var block *btcjson.GetBlockVerboseResult
	decoder := gob.NewDecoder(bytes.NewReader(x))
	decoder.Decode(&block)

	json.NewEncoder(w).Encode(&block)

	// x := database.ViewBlock("c65aabe1578da37945118ff6078792315499dd0dd6712f76a5f387126799d9b1")
	//
	// var block *btcjson.GetBlockVerboseResult
	// decoder := gob.NewDecoder(bytes.NewReader(x))
	// decoder.Decode(&block)
	// x := database.FetchTransactionHashByBlockhash("c65aabe1578da37945118ff6078792315499dd0dd6712f76a5f387126799d9b1")

	// var block *btcjson.TxRawResult
	// decoder := gob.NewDecoder(bytes.NewReader(x))
	// decoder.Decode(&block)
	//

	// var block *btcjson.TxRawDecodeResult
	// //fmt.Println(btcjson.DecodeRawTransactionCmd{x})
	// decoder := gob.NewDecoder(bytes.NewReader(x))
	// decoder.Decode(&block)

	// json.NewEncoder(w).Encode(string(x))
	// fmt.Println(btcjson.NewDecodeRawTransactionCmd(string(x)))
}
