package daos

import (
	"time"
)

// TODO: implement reading from config
var dbHosts = []string{"localhost"}

type dbName string

const (
	// Coin <~> Database
	Via dbName = "viacoin"
	Btc dbName = "Btc"
	Ltc dbName = "litecoin"

	// timeouts
	timeout10 = 10 * time.Second
	timeout60 = 60 * time.Second

	// Collection names
	txs    = "Txs"
	blocks = "Blocks"
)

// // Database dialing info
// var (
// 	viaDialInfo = &mgo.DialInfo{
// 		Addrs:    dbHosts,
// 		Timeout:  timeout10,
// 		Database: Via,
// 	}
//
// 	btcDialInfo = &mgo.DialInfo{
// 		Addrs:    dbHosts,
// 		Timeout:  timeout10,
// 		Database: Btc,
// 	}
//
// 	ltcDialInfo = &mgo.DialInfo{
// 		Addrs:    dbHosts,
// 		Timeout:  timeout10,
// 		Database: Ltc,
// 	}
// )
