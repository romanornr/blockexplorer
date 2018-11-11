package daos

import (
	"time"
)

// TODO: implement reading from config
var dbHosts = []string{"localhost"}

type dbName string

const (
	// Coin <~> Database
	Via     dbName = "viacoin"
	ViaTest dbName = "viatest"

	Btc     dbName = "Btc"
	BtcTest dbName = "btctest"

	Ltc     dbName = "litecoin"
	LtcTest dbName = "ltctest"

	// timeouts
	timeout10 = 10 * time.Second
	timeout60 = 60 * time.Second

	// Collection names
	txs    = "Txs"
	blocks = "Blocks"
	addr   = "Addresses"
)
