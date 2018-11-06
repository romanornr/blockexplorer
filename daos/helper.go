package daos

import (
	"time"
)

// TODO: implement reading from config
var dbHosts = []string{"localhost"}

type dbConn string

const (
	// Coin <~> Database
	Via dbConn = "viacoin"
	Btc dbConn = "Btc"
	Ltc dbConn = "litecoin"

	// timeouts
	timeout10 = 10 * time.Second
	timeout60 = 60 * time.Second
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
