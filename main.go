package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded ! Please check the config folder")
	}

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())
	fmt.Printf("using %s:%d\n", viper.Get("server.ip"), viper.Get("server.port"))
}
