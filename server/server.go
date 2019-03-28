package server

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("website/*"))

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())
	fmt.Printf("Listening on %s:%d\n", viper.GetString("server.ip"), viper.Get("server.port"))
}

func Start() {

	//go zeroMQ.BlockNotify()   uncomment this to get new blocks added. Commented out due to development now.

	port := ":" + viper.GetString("server.port")
	addr := flag.String("addr", port, "http service address")

	flag.Parse()
	fmt.Println("Server started...")

	router := createRouter()
	err := http.ListenAndServe(*addr, router)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
