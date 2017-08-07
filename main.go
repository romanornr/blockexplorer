package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/btcsuite/btcrpcclient"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
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

func client() *btcrpcclient.Client {

	// Connect to local bitcoin/altcoin core RPC server using HTTP POST mode.
	connCfg := &btcrpcclient.ConnConfig{
		Host:         viper.GetString("rpc.ip") + ":" + viper.GetString("rpc.port"), //127.0.0.1:8332
		User:         viper.GetString("rpc.username"),
		Pass:         viper.GetString("rpc.password"),
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := btcrpcclient.New(connCfg, nil)
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
}

func Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	coin := viper.Get("coin.name")
	err := tpl.ExecuteTemplate(w, "index.html", coin)
	if err != nil {
		log.Println("errror")
	}
}

func GetDifficulty(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	info, err := client().GetDifficulty()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(info)
}

// GetLatestBlocks gets 10 of the latest blocks
func GetLatestBlocks(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	blockCount, err := client().GetBlockCount()
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < 10; i++ {
		prevBlock := blockCount - int64(i)
		hash, err := client().GetBlockHash(prevBlock)
		if err != nil {
			log.Println(err)
		}

		test2, err := client().GetBlock(hash)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%d:%s\n", prevBlock, hash)

		data, err := json.Marshal(test2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", data)
		json.NewEncoder(w).Encode(data)
	}

}
