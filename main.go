package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/btcsuite/btcrpcclient"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

var tpl *template.Template

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
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

	router := httprouter.New()
	router.GET("/", Index)
	http.ListenAndServe(viper.GetString("server.ip")+":"+viper.GetString("server.port"), router) //example: 127.0.0.1:8080
}

func Index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	coin := viper.Get("coin.name")
	err := tpl.ExecuteTemplate(w, "index.html", coin)
	if err != nil {
		log.Println("errror")
	}
}
