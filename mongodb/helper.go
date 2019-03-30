package mongodb

import (
	"github.com/globalsign/mgo"
	"time"
)

// TODO: implement reading from config

//var x := viper.GetString("mongodb.ip")
var dbHosts = []string{"localhost"}

type dbName string

const (
	// Coin <~> Database
	VIA  dbName = "viacoin"
	BTC  dbName = "bitcoin"
	LTC  dbName = "litecoin"
	TEST dbName = "explorer-test"

	// timeouts
	timeout10 = 10 * time.Second
	timeout60 = 60 * time.Second

	// Collection names
	TRANSACTIONS = "Transactions"
	BLOCKS       = "Blocks"
	ADDRESSINFO  = "AddressInfo"
)

var (
	viaDialInfo = &mgo.DialInfo{
		Addrs:    dbHosts,
		Timeout:  timeout10,
		Database: string(VIA),
	}
)
