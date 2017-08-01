package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

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
}

func main() {

	mux := httprouter.New()
	mux.GET("/", index)
	http.ListenAndServe(":8000", mux)

	//fmt.Printf("using %s:%d\n", viper.Get("server.ip"), viper.Get("server.port"))
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	coin := viper.Get("coin.name")
	err := tpl.ExecuteTemplate(w, "index.html", coin)
	if err != nil {
		log.Println("errror")
	}
}
