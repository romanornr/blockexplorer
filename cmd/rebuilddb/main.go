package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"log"
	"github.com/spf13/viper"
	"github.com/romanornr/cyberchain/blockdata"
)

var redisClient *redis.Client

func init(){
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file loaded ! Please check the config folder")
	}

	fmt.Printf("Reading configuration from %s\n", viper.ConfigFileUsed())
}

func main(){
	fmt.Println(blockdata.GetBlockHash(0))

	//e := blockdata.GetBlockHash(0).String()
	err := redisClient.Set("0", blockdata.GetBlockHash(0).String(), 0).Err()
	if err != nil{
		panic(err)
	}
	//fmt.Println(blockdata.GetDifficulty())
}
