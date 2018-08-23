package server

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"github.com/viacoin/viad/rpcclient"
)

var tpl *template.Template

func client() *rpcclient.Client {
	// Connect to local bitcoin/altcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         viper.GetString("rpc.ip") + ":" + viper.GetString("rpc.port"), // 127.0.0.1:8332
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

	return client
}

func init() {
	tpl = template.Must(template.ParseGlob("website/*"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded! Please check the config folder")
	}

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())
	fmt.Printf("Listening on %s:%d\n", viper.GetString("server.ip"), viper.Get("server.port"))
}

func Start() {
	port := ":" + viper.GetString("server.port")
	addr := flag.String("addr", port, "http service address")

	flag.Parse()

	router := createRouter()
	err := http.ListenAndServe(*addr, router)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
