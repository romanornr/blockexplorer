package client

import (
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/spf13/viper"
	"log"
)

var instance *rpcclient.Client

func GetInstance() *rpcclient.Client {
	if instance == nil {
		connCfg := LoadConfig()
		var err error
		instance, err = rpcclient.New(connCfg, nil)
		if err != nil {
			log.Fatal(err)
			instance.Shutdown()
		}
	}
	return instance
}

func LoadConfig() *rpcclient.ConnConfig{
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded ! Please check the config folder")
	}

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())

	connCfg := &rpcclient.ConnConfig{
		Host:         viper.GetString("rpc.ip") + ":" + viper.GetString("rpc.port"), //127.0.0.1:8332
		User:         viper.GetString("rpc.username"),
		Pass:         viper.GetString("rpc.password"),
		HTTPPostMode: true, // Viacoin core only supports HTTP POST mode
		DisableTLS:   true, // Viacoin core does not provide TLS by default
	}
	return connCfg
}
