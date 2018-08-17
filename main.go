package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"github.com/btcsuite/btcd/rpcclient"
	_"encoding/gob"
	_"bytes"
	"github.com/romanornr/cyberchain/database"
	"encoding/gob"
	"bytes"
)

var tpl *template.Template

//Handlers run concurrently but maps aren't thread-safe
//a Mutex is used to ensure only 1 goroutine can update data
type store struct {
	data map[string]string
	m    sync.RWMutex
}

//create the datastore
var (
	s = store{
		data: map[string]string{},
		m:    sync.RWMutex{},
	}
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

func init() {
	tpl = template.Must(template.ParseGlob("website/*"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded ! Please check the config folder")
	}

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())
	fmt.Printf("Webserving starting using %s:%d\n", viper.GetString("server.ip"), viper.Get("server.port"))
}

func main() {

	port := ":" + viper.GetString("server.port")
	addr := flag.String("addr", port, "http service address")

	flag.Parse()
	network := viper.GetString("coin.symbol")

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/"+network+"/getdifficulty", GetDifficulty)
	router.GET("/api/"+network+"/blocks", GetLatestBlocks)
	router.GET("/api/"+network+"/block",
		GetBlock)

	fileServer := http.FileServer(http.Dir("static"))
	router.GET("/static/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		fileServer.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(*addr, router)
	if err != nil {
		log.Fatal("ListenAdnServe:", err)
	}

	defer client().Shutdown()
}

func Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	coin := viper.Get("coin.name")
	err := tpl.ExecuteTemplate(w, "index.html", coin)
	if err != nil {
		log.Println("errror")
	}
}

func GetDifficulty(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	difficulty, err := client().GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(difficulty)
}

// GetxLatestBlocks gets x (int) latest blocks
func GetLatestBlocks(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	blockCount, err := client().GetBlockCount() //get the latest blocks
	if err != nil {
		log.Println(err)
	}

	var blocks []*btcjson.GetBlockVerboseResult

	// blockheight - 1 in the loop. Get the blockhash from the height
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

func GetBlock(w http.ResponseWriter, req *http.Request, _ httprouter.Params){

	x := database.ViewBlock("c65aabe1578da37945118ff6078792315499dd0dd6712f76a5f387126799d9b1")

	var block *btcjson.GetBlockVerboseResult
	decoder := gob.NewDecoder(bytes.NewReader(x))
	decoder.Decode(&block)

	json.NewEncoder(w).Encode(&block)
	//x := database.ViewBlock("c65aabe1578da37945118ff6078792315499dd0dd6712f76a5f387126799d9b1")
	//
	//var block *btcjson.GetBlockVerboseResult
	//decoder := gob.NewDecoder(bytes.NewReader(x))
	//decoder.Decode(&block)
	//x := database.FetchTransactionHashByBlockhash("c65aabe1578da37945118ff6078792315499dd0dd6712f76a5f387126799d9b1")

	//var block *btcjson.TxRawResult
	//decoder := gob.NewDecoder(bytes.NewReader(x))
	//decoder.Decode(&block)
	//

	//var block *btcjson.TxRawDecodeResult
	////fmt.Println(btcjson.DecodeRawTransactionCmd{x})
	//decoder := gob.NewDecoder(bytes.NewReader(x))
	//decoder.Decode(&block)

	//json.NewEncoder(w).Encode(string(x))
	//fmt.Println(btcjson.NewDecodeRawTransactionCmd(string(x)))
}
